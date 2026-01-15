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
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	v1 "github.com/odysseia-greek/apologia/aspasia/gen/go/v1"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
)

func (g *GathererServiceImpl) Search(ctx context.Context, request *v1.ExtendedSearch) (*v1.ExtendedSearchResponse, error) {
	requestId := CurrentRequestID(ctx, config.DefaultTracingName, service.HeaderKey)

	analyseResult := &v1.ExtendedSearchResponse{}

	cacheItem, _ := g.Archytas.Read(request.Word)
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
		dictionaryResults, err = g.gatherSimilarWords(egCtx, request.Word)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	analyseResult.FoundInText = textResults
	analyseResult.SimilarWords = dictionaryResults

	analyseResultJson, _ := json.Marshal(analyseResult)

	standardDuration := time.Minute * 10
	err := g.Archytas.SetWithTTL(request.Word, string(analyseResultJson), standardDuration)
	if err != nil {
		logging.Error(err.Error())
	}

	return analyseResult, nil
}

func (g *GathererServiceImpl) gatherSimilarWords(ctx context.Context, word string) ([]*v1.SimilarWords, error) {
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

	if rid, _ := ctx.Value(config.HeaderKey).(string); rid != "" {
		req.Header.Set(config.HeaderKey, rid)
	}
	if sid, _ := ctx.Value(config.SessionIdKey).(string); sid != "" {
		req.Header.Set(config.SessionIdKey, sid)
	}

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

	return similarWords, nil
}

func (g *GathererServiceImpl) gatherTexts(ctx context.Context, word, requestId string) (*v1.AnalyzeTextResponse, error) {
	var analyseResult *v1.AnalyzeTextResponse

	herodotosSpan := &arv1.ObserveRequest{
		Kind: &arv1.ObserveRequest_Action{Action: &arv1.ObserveAction{
			Action: "analyseText",
			Status: fmt.Sprintf("querying Herodotos for word: %s", word),
		}},
	}

	comedy.ServiceToServiceSpanWithCtx(ctx, herodotosSpan, g.Streamer)

	r := models.AnalyzeTextRequest{Rootword: word}
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	foundInText, err := g.Client.Herodotos().Analyze(jsonBody, requestId)

	if foundInText != nil {
		var source models.AnalyzeTextResponse
		defer foundInText.Body.Close()
		err = json.NewDecoder(foundInText.Body).Decode(&source)
		if err != nil {
			logging.Error(fmt.Sprintf("error while decoding: %s", err.Error()))
		}

		analyseResult = &v1.AnalyzeTextResponse{
			Rootword:     source.Rootword,
			PartOfSpeech: source.PartOfSpeech,
			Conjugations: []*v1.Conjugations{},
			Texts:        []*v1.AnalyzeResult{},
		}

		for _, text := range source.Results {
			parsedText := &v1.AnalyzeResult{
				ReferenceLink: text.ReferenceLink,
				Author:        text.Author,
				Book:          text.Book,
				Reference:     text.Reference,
				Text: &v1.Rhema{
					Greek:        text.Text.Greek,
					Translations: text.Text.Translations,
					Section:      text.Text.Section,
				},
			}
			analyseResult.Texts = append(analyseResult.Texts, parsedText)
		}

		for _, conjugation := range source.Conjugations {
			parsedConjugation := &v1.Conjugations{
				Word: conjugation.Word,
				Rule: conjugation.Rule,
			}

			analyseResult.Conjugations = append(analyseResult.Conjugations, parsedConjugation)
		}

		logging.Debug(fmt.Sprintf("found in herodotos: %s number of results: %d", word, len(analyseResult.Texts)))
	}

	return analyseResult, nil
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
