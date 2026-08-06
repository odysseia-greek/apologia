package rhetorike

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	dionysiosv1 "github.com/odysseia-greek/alexandreia/dionysios/gen/go/v1"
	v1 "github.com/odysseia-greek/apologia/aspasia/gen/go/v1"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
)

func (g *GathererServiceImpl) Search(ctx context.Context, request *v1.ExtendedSearch) (*v1.ExtendedSearchResponse, error) {
	requestId := CurrentRequestID(ctx, config.DefaultTracingName, service.HeaderKey)

	cleanWord := parseWordIntoParts(request.Word)
	analyseResult := &v1.ExtendedSearchResponse{}

	request.Word = cleanWord

	cacheItem, _ := g.Archytas.Get(request.Word)
	if cacheItem != nil {
		err := json.Unmarshal(cacheItem, &analyseResult)
		if err != nil {
			return nil, err
		}

		logging.Debug(fmt.Sprintf("found in cache: %s number of results: %d", request.Word, len(analyseResult.FoundInText.Texts)))

		go comedy.CacheSpan(string(cacheItem), request.Word, ctx, g.Streamer)
		return analyseResult, nil
	}

	var textResults *v1.AnalyzeTextResponse
	var dictionaryResults []*v1.SimilarWords

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		var err error
		textResults, err = g.gatherTexts(egCtx, request.Word, requestId)
		return err
	})

	eg.Go(func() error {
		var err error
		dictionaryResults, err = g.gatherSimilarWords(egCtx, request.Word, requestId)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	analyseResult.FoundInText = textResults
	analyseResult.SimilarWords = dictionaryResults

	analyseResultJson, _ := json.Marshal(analyseResult)

	standardDuration := time.Minute * 10
	err := g.Archytas.SetBytesWithTTL(request.Word, analyseResultJson, standardDuration)
	if err != nil {
		logging.Error(err.Error())
	}

	return analyseResult, nil
}

func (g *GathererServiceImpl) gatherSimilarWords(ctx context.Context, word, requestId string) ([]*v1.SimilarWords, error) {
	var similarWords []*v1.SimilarWords

	// ---- tracing: start action span
	alexandrosSpan := &arv1.ObserveRequest{
		Kind: &arv1.ObserveRequest_Action{Action: &arv1.ObserveAction{
			Action: "gatherSimilarWords",
			Status: fmt.Sprintf("querying Alexandros fuzzy for word: %s", word),
		}},
	}
	// This should set TraceId / ParentSpanId / SpanId using ctx (your helper)
	comedy.ServiceToServiceSpanWithCtx(ctx, alexandrosSpan, g.Streamer)

	start := time.Now()
	defer func() {
		// close span as action (or TraceHopStop if you prefer)
		if g.Streamer != nil {
			_ = g.Streamer.Send(&arv1.ObserveRequest{
				TraceId:      alexandrosSpan.TraceId,
				ParentSpanId: alexandrosSpan.SpanId,
				SpanId:       comedy.GenerateSpanID(),
				Kind: &arv1.ObserveRequest_Action{
					Action: &arv1.ObserveAction{
						Action: "CloseSpan",
						Status: "alexandros fuzzy finished",
						TookMs: time.Since(start).Milliseconds(),
					},
				},
			})
		}
	}()

	// ---- build GraphQL request
	body := gqlReq{
		Query: fuzzyMiniQuery,
		Variables: map[string]any{
			"input": map[string]any{
				"word": word,
				"size": 20, // ask more than we need; we’ll trim to 5 after filtering
			},
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// ---- make HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.AlexandrosAddress, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(config.HeaderKey, requestId)

	client := g.GraphqlClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("alexandros fuzzy http %d: %s", res.StatusCode, string(raw))
	}

	var parsed alexandrosFuzzyResp
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	if len(parsed.Errors) > 0 {
		return nil, fmt.Errorf("alexandros fuzzy gql error: %s", parsed.Errors[0].Message)
	}

	// ---- filter / dedupe / trim to 5
	targetNorm := g.normalizeGreekWithArticle(word)
	seen := make(map[string]struct{}, 16)

	for _, r := range parsed.Data.Fuzzy.Results {
		if len(similarWords) >= 5 {
			break
		}
		if r.Headword == "" {
			continue
		}

		headwordNorm := g.normalizeGreekWithArticle(r.Headword)

		// skip identical to query word
		if headwordNorm == targetNorm {
			continue
		}

		// skip duplicates (handles "λόγος" vs "ὁ λόγος" etc)
		if _, ok := seen[headwordNorm]; ok {
			continue
		}
		seen[headwordNorm] = struct{}{}

		sw := &v1.SimilarWords{
			Greek: r.Headword,
		}

		// pick english quick gloss if present
		for _, qg := range r.QuickGlosses {
			if qg.Language == "en" {
				sw.English = qg.Gloss
				break
			}
		}

		similarWords = append(similarWords, sw)
	}

	logging.Debug(fmt.Sprintf("found in alexandros: %s number of results: %d", word, len(similarWords)))

	return similarWords, nil
}

func (g *GathererServiceImpl) gatherTexts(ctx context.Context, word, requestId string) (*v1.AnalyzeTextResponse, error) {
	dionysiosSpan := &arv1.ObserveRequest{
		Kind: &arv1.ObserveRequest_Action{Action: &arv1.ObserveAction{
			Action: "researchText",
			Status: fmt.Sprintf("querying Dionysios for word: %s", word),
		}},
	}

	comedy.ServiceToServiceSpanWithCtx(ctx, dionysiosSpan, g.Streamer)

	if g.Dionysios == nil {
		return nil, fmt.Errorf("Dionysios research client is not configured")
	}

	researchCtx, cancel := createRequestHeader(ctx, requestId, "")
	defer cancel()

	source, err := g.Dionysios.Research(researchCtx, &dionysiosv1.ResearchRequest{
		Rootword: word,
		Limit:    5,
	})
	if err != nil {
		return nil, fmt.Errorf("Dionysios research failed: %w", err)
	}

	analyseResult := mapDionysiosResearch(source)
	logging.Debug(fmt.Sprintf("found in Dionysios: %s number of results: %d", word, len(analyseResult.Texts)))
	return analyseResult, nil
}

func mapDionysiosResearch(source *dionysiosv1.ResearchResponse) *v1.AnalyzeTextResponse {
	result := &v1.AnalyzeTextResponse{
		Conjugations: []*v1.Conjugations{},
		Texts:        []*v1.AnalyzeResult{},
	}
	if source == nil {
		return result
	}

	result.Rootword = source.GetRootword()
	result.PartOfSpeech = source.GetPartOfSpeech()
	for _, conjugation := range source.GetConjugations() {
		result.Conjugations = append(result.Conjugations, &v1.Conjugations{
			Word: conjugation.GetWord(),
			Rule: conjugation.GetRule(),
		})
	}
	for _, text := range source.GetResults() {
		mapped := &v1.AnalyzeResult{
			ReferenceLink: text.GetReferenceLink(),
			Author:        text.GetAuthor(),
			Book:          text.GetBook(),
			Reference:     text.GetReference(),
		}
		if text.GetText() != nil {
			mapped.Text = &v1.Rhema{
				Greek:        text.GetText().GetGreek(),
				Translations: text.GetText().GetTranslations(),
				Section:      text.GetText().GetSection(),
			}
		}
		result.Texts = append(result.Texts, mapped)
	}

	return result
}

func (g *GathererServiceImpl) normalizeGreekWithArticle(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	articles := []string{"ὁ", "ἡ", "τό"}
	for _, a := range articles {
		prefix := a + " "
		if strings.HasPrefix(s, prefix) {
			s = strings.TrimSpace(strings.TrimPrefix(s, prefix))
			break
		}
	}

	return s
}

func createRequestHeader(ctx context.Context, requestID, sessionId string) (context.Context, context.CancelFunc) {
	requestCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	md := metadata.New(map[string]string{
		config.HeaderKey:    requestID,
		config.SessionIdKey: sessionId,
	})
	return metadata.NewOutgoingContext(requestCtx, md), cancel
}

func CurrentRequestID(ctx context.Context, ctxKey any, headerKey string) string {
	if v := ctx.Value(ctxKey); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get(headerKey); len(vals) > 0 {
			return vals[0]
		}
	}
	return ""
}

func parseWordIntoParts(word string) string {
	word = strings.TrimSpace(word)

	// Split on whitespace
	parts := strings.Fields(word)
	if len(parts) <= 1 {
		return word
	}

	// Known Greek articles (with and without accents)
	articles := map[string]struct{}{
		"ὁ": {}, "ἡ": {}, "τό": {}, "τὸ": {},
		"οἱ": {}, "αἱ": {}, "τά": {}, "τὰ": {},
	}

	// If first part is an article, return the rest
	if _, ok := articles[parts[0]]; ok {
		return strings.Join(parts[1:], " ")
	}

	// Fallback: return last part (safer than first)
	return parts[len(parts)-1]
}
