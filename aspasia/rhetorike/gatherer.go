package rhetorike

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	v1 "github.com/odysseia-greek/apologia/aspasia/gen/go/v1"
	"github.com/odysseia-greek/apologia/diotima/theoria"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	antigonosv1 "github.com/odysseia-greek/makedonia/antigonos/gen/go/v1"
	"github.com/odysseia-greek/makedonia/antigonos/monophthalmus"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
)

func (g *GathererServiceImpl) Search(ctx context.Context, request *v1.ExtendedSearch) (*v1.ExtendedSearchResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var requestId string
	if ok {
		headerValue := md.Get(service.HeaderKey)
		if len(headerValue) > 0 {
			requestId = headerValue[0]
		}
	}

	analyseResult := &v1.ExtendedSearchResponse{}

	cacheItem, _ := g.Archytas.Read(request.Word)
	if cacheItem != nil {
		err := json.Unmarshal(cacheItem, &analyseResult)
		if err != nil {
			return nil, err
		}

		logging.Debug(fmt.Sprintf("found in cache: %s number of results: %d", request.Word, len(analyseResult.FoundInText.Texts)))

		go theoria.CacheSpan(string(cacheItem), request.Word, ctx)
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

	antigonosSpan := &pb.ParabasisRequest{
		RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
			Action: "analyseText",
			Status: fmt.Sprintf("querying Antigonos for word: %s", word),
		}},
	}

	fuzzyClientCtx, cancel := g.createRequestHeader(requestId, "")
	defer cancel()

	theoria.ServiceToServiceSpan(antigonosSpan, ctx)

	request := &koinos.SearchQuery{
		Word:            word,
		Language:        koinos.Language_LANG_GREEK,
		NumberOfResults: 10, // we search for 10 words so we can return at least 5 results
	}
	startTime := time.Now()

	var grpcResponse *antigonosv1.SearchResponse

	err := g.FuzzyClient.CallWithReconnect(func(client *monophthalmus.FuzzyClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Search(fuzzyClientCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	endTime := time.Since(startTime)
	antigonosSpan = &pb.ParabasisRequest{
		RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
			Action: "fuzzySearch",
			Took:   fmt.Sprintf("%v", endTime),
			Status: "querying Antigonos returned success",
		}},
	}
	theoria.ServiceToServiceSpan(antigonosSpan, ctx)

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
	herodotosSpan := &pb.ParabasisRequest{
		RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
			Action: "analyseText",
			Status: fmt.Sprintf("querying Herodotos for word: %s", word),
		}},
	}

	theoria.ServiceToServiceSpan(herodotosSpan, ctx)

	startTime := time.Now()
	r := models.AnalyzeTextRequest{Rootword: word}
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	foundInText, err := g.Client.Herodotos().Analyze(jsonBody, requestId)
	endTime := time.Since(startTime)

	if foundInText != nil {
		var source models.AnalyzeTextResponse
		defer foundInText.Body.Close()
		err = json.NewDecoder(foundInText.Body).Decode(&source)
		if err != nil {
			logging.Error(fmt.Sprintf("error while decoding: %s", err.Error()))
		}

		herodotosSpan = &pb.ParabasisRequest{
			RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
				Action: "analyseText",
				Took:   fmt.Sprintf("%v", endTime),
				Status: fmt.Sprintf("querying Herodotos returned: %d", foundInText.StatusCode),
			}},
		}
		theoria.ServiceToServiceSpan(herodotosSpan, ctx)

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
