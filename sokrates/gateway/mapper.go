package gateway

import (
	"github.com/odysseia-greek/agora/plato/models"
	pbartrippos "github.com/odysseia-greek/apologia/aristippos/proto"
)

func mapComprehensiveResponse(grpcResp *pbartrippos.ComprehensiveResponse) *models.ComprehensiveResponse {
	if grpcResp == nil {
		return nil
	}

	mappedResponse := &models.ComprehensiveResponse{
		Correct:  grpcResp.Correct,
		QuizWord: grpcResp.QuizWord,
	}

	if grpcResp.FoundInText != nil {
		mappedResponse.FoundInText = models.AnalyzeTextResponse{
			Rootword:     grpcResp.FoundInText.Rootword,
			PartOfSpeech: grpcResp.FoundInText.PartOfSpeech,
			Conjugations: mapConjugations(grpcResp.FoundInText.Conjugations),
			Results:      mapAnalyzeResults(grpcResp.FoundInText.Texts),
		}
	}

	for _, word := range grpcResp.SimilarWords {
		mappedResponse.SimilarWords = append(mappedResponse.SimilarWords, models.Meros{
			Greek:      word.Greek,
			English:    word.English,
			Dutch:      word.Dutch,
			LinkedWord: word.LinkedWord,
			Original:   word.Original,
		})
	}

	return mappedResponse
}

func mapConjugations(grpcConj []*pbartrippos.Conjugations) []models.Conjugations {
	var result []models.Conjugations
	for _, conj := range grpcConj {
		result = append(result, models.Conjugations{
			Word: conj.Word,
			Rule: conj.Rule,
		})
	}
	return result
}

func mapAnalyzeResults(grpcResults []*pbartrippos.AnalyzeResult) []models.AnalyzeResult {
	var result []models.AnalyzeResult
	for _, res := range grpcResults {
		result = append(result, models.AnalyzeResult{
			ReferenceLink: res.ReferenceLink,
			Author:        res.Author,
			Book:          res.Book,
			Reference:     res.Reference,
			Text: models.Rhema{
				Greek:        res.Text.Greek,
				Translations: res.Text.Translations,
				Section:      res.Text.Section,
			},
		})
	}
	return result
}
