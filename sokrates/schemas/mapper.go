package schemas

import (
	"github.com/graphql-go/graphql"
	pbartrippos "github.com/odysseia-greek/apologia/aristippos/proto"
)

func buildQuizRequest(p graphql.ResolveParams) *pbartrippos.CreationRequest {
	theme, _ := p.Args["theme"].(string)
	segment, _ := p.Args["segment"].(string)
	order, _ := p.Args["order"].(string)
	excludeWords, _ := p.Args["excludeWords"].([]interface{})

	excludeWordsStr := make([]string, len(excludeWords))
	for i, word := range excludeWords {
		excludeWordsStr[i], _ = word.(string)
	}

	set, isOK := p.Args["set"].(string)
	if !isOK {
		return nil // should return an error if necessary
	}

	return &pbartrippos.CreationRequest{
		Theme:        theme,
		Set:          set,
		Segment:      segment,
		Order:        order,
		ExcludeWords: excludeWordsStr,
	}
}

func buildAnswerRequest(p graphql.ResolveParams) *pbartrippos.AnswerRequest {
	theme, _ := p.Args["theme"].(string)
	segment, _ := p.Args["segment"].(string)
	quizWord, _ := p.Args["quizWord"].(string)
	answer, _ := p.Args["answer"].(string)
	comprehensive, _ := p.Args["comprehensive"].(bool)

	set, isOK := p.Args["set"].(string)
	if !isOK {
		return nil // should return an error if necessary
	}
	quizType, isOK := p.Args["quizType"].(string)
	if !isOK {
		return nil // should return an error if necessary
	}

	return &pbartrippos.AnswerRequest{
		Theme:         theme,
		Set:           set,
		Segment:       segment,
		QuizType:      quizType,
		Comprehensive: comprehensive,
		Answer:        answer,
		QuizWord:      quizWord,
	}
}
