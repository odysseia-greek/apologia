package multiplechoice

import (
	pbkritias "github.com/odysseia-greek/apologia/kritias/proto"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func MapComprehensiveResponse(grpcResp *pbkritias.ComprehensiveResponse) *model.ComprehensiveResponse {
	if grpcResp == nil {
		return nil
	}

	mappedResponse := &model.ComprehensiveResponse{
		Correct:  &grpcResp.Correct,
		QuizWord: &grpcResp.QuizWord,
	}

	if grpcResp.FoundInText != nil {
		mappedResponse.FoundInText = &model.AnalyzeTextResponse{
			Rootword:     &grpcResp.FoundInText.Rootword,
			Conjugations: mapConjugations(grpcResp.FoundInText.Conjugations),
			Texts:        mapAnalyzeResults(grpcResp.FoundInText.Texts),
		}
	}

	for _, word := range grpcResp.SimilarWords {
		mappedResponse.SimilarWords = append(mappedResponse.SimilarWords, &model.Hit{
			Greek:      &word.Greek,
			English:    &word.English,
			Dutch:      &word.Dutch,
			LinkedWord: &word.LinkedWord,
			Original:   &word.Original,
		})
	}

	for _, progress := range grpcResp.Progress {
		mappedResponse.Progress = append(mappedResponse.Progress, &model.ProgressEntry{
			Greek:          &progress.Greek,
			Translation:    &progress.Translation,
			PlayCount:      &progress.PlayCount,
			CorrectCount:   &progress.CorrectCount,
			IncorrectCount: &progress.IncorrectCount,
			LastPlayed:     &progress.LastPlayed,
		})
	}

	return mappedResponse
}

func mapConjugations(grpcConj []*pbkritias.Conjugations) []*model.ConjugationResponse {
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

func mapAnalyzeResults(grpcResults []*pbkritias.AnalyzeResult) []*model.AnalyzeResult {
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

// Helper function to convert []string to []*string
func convertStringSliceToPointer(strings []string) []*string {
	var ptrSlice []*string
	for _, s := range strings {
		ptrSlice = append(ptrSlice, &s)
	}
	return ptrSlice
}
