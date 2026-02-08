package gateway

import (
	"context"

	v1 "github.com/odysseia-greek/apologia/aspasia/gen/go/v1"
	"github.com/odysseia-greek/apologia/aspasia/rhetorike"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) GatherComprehensiveResponse(ctx context.Context, word string) (*model.ComprehensiveResponse, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.ExtendedSearchResponse

	request := &v1.ExtendedSearch{Word: word}

	err := s.GathererClient.CallWithReconnect(func(client *rhetorike.GathererClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Search(outCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	return mapComprehensiveResponse(grpcResponse), nil
}

func mapComprehensiveResponse(grpcResp *v1.ExtendedSearchResponse) *model.ComprehensiveResponse {
	if grpcResp == nil {
		return nil
	}

	mappedResponse := &model.ComprehensiveResponse{}

	if grpcResp.FoundInText != nil {
		mappedResponse.FoundInText = &model.AnalyzeTextResponse{
			Rootword:     &grpcResp.FoundInText.Rootword,
			Conjugations: mapConjugations(grpcResp.FoundInText.Conjugations),
			Texts:        mapAnalyzeResults(grpcResp.FoundInText.Texts),
		}
	}

	for _, word := range grpcResp.SimilarWords {
		mappedResponse.SimilarWords = append(mappedResponse.SimilarWords, &model.Hit{
			Greek:   &word.Greek,
			English: &word.English,
		})
	}

	return mappedResponse
}

func mapConjugations(grpcConj []*v1.Conjugations) []*model.ConjugationResponse {
	if grpcConj == nil {
		return nil
	}

	var result []*model.ConjugationResponse
	for _, conj := range grpcConj {
		result = append(result, &model.ConjugationResponse{
			Word: &conj.Word,
			Rule: &conj.Rule,
		})
	}
	return result
}

func mapAnalyzeResults(grpcResults []*v1.AnalyzeResult) []*model.AnalyzeResult {
	if grpcResults == nil {
		return nil
	}

	var result []*model.AnalyzeResult
	for _, res := range grpcResults {
		result = append(result, &model.AnalyzeResult{
			ReferenceLink: &res.ReferenceLink,
			Author:        &res.Author,
			Book:          &res.Book,
			Reference:     &res.Reference,
			Text: &model.Rhema{
				Greek:        &res.Text.Greek,
				Translations: convertStringSliceToPointer(res.Text.Translations),
				Section:      &res.Text.Section,
			},
		})
	}
	return result
}
