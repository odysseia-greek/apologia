package rhetorike

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	v1 "github.com/odysseia-greek/apologia/aspasia/gen/go/v1"
	koinosv1 "github.com/odysseia-greek/apologia/diotima/gen/go/koinos/v1"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	antigonosv1 "github.com/odysseia-greek/makedonia/antigonos/gen/go/v1"
	"github.com/odysseia-greek/makedonia/antigonos/monophthalmus"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
)

func (g *GathererServiceImpl) Health(context.Context, *koinosv1.HealthRequest) (*koinosv1.HealthResponse, error) {
	return &koinosv1.HealthResponse{
		Healthy:        true,
		Time:           time.Now().String(),
		DatabaseHealth: nil,
		Version:        g.Version,
	}, nil
}

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
	err := g.Archytas.SetWithTTL(request.Word, string(analyseResultJson), standardDuration)
	if err != nil {
		logging.Error(err.Error())
	}

	return analyseResult, nil
}

func (g *GathererServiceImpl) gatherSimilarWords(
	ctx context.Context,
	word, requestId string,
) ([]*v1.SimilarWords, error) {
	var similarWords []*v1.SimilarWords

	antigonosSpan := &arv1.ObserveRequest{
		Kind: &arv1.ObserveRequest_Action{Action: &arv1.ObserveAction{
			Action: "analyseText",
			Status: fmt.Sprintf("querying Antigonos for word: %s", word),
		}},
	}

	fuzzyClientCtx, cancel := createRequestHeader(ctx, requestId, "")
	defer cancel()

	comedy.ServiceToServiceSpanWithCtx(ctx, antigonosSpan, g.Streamer)

	request := &koinos.SearchQuery{
		Word:            word,
		Language:        koinos.Language_LANG_GREEK,
		NumberOfResults: 10, // we search for 10 words so we can return at least 5 results
	}

	var grpcResponse *antigonosv1.SearchResponse

	err := g.FuzzyClient.CallWithReconnect(func(client *monophthalmus.FuzzyClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Search(fuzzyClientCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	// Precompute normalized form of the requested word
	targetNorm := g.normalizeGreekWithArticle(word)

	for _, result := range grpcResponse.Results {
		if len(similarWords) >= 5 {
			break
		}

		// Normalize headword and skip if it's effectively the same as the input.
		headwordNorm := g.normalizeGreekWithArticle(result.Headword)
		if headwordNorm == targetNorm {
			continue
		}

		// avoid duplicates if Antigonos returns both "λόγος" and "ὁ λόγος").
		duplicate := false
		for _, existing := range similarWords {
			if g.normalizeGreekWithArticle(existing.Greek) == headwordNorm {
				duplicate = true
				break
			}
		}
		if duplicate {
			continue
		}

		similarWord := v1.SimilarWords{
			Greek: result.Headword,
		}

		for _, gloss := range result.QuickGlosses {
			if gloss.Language == "en" {
				similarWord.English = gloss.Gloss
				break
			}
		}

		similarWords = append(similarWords, &similarWord)
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
