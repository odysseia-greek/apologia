package gateway

import (
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
	pbxenofon "github.com/odysseia-greek/apologia/xenofon/proto"
)

func (s *SokratesHandler) CreateAuthorBasedQuiz(request *pbxenofon.CreationRequest, requestID, sessionId string) (*model.AuthorBasedResponse, error) {
	authorBasedCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.AuthorBasedClient.Question(authorBasedCtx, request)
	if err != nil {
		return nil, err
	}

	quizResponse := &model.AuthorBasedResponse{
		FullSentence: &grpcResponse.FullSentence,
		Translation:  &grpcResponse.Translation,
		Reference:    &grpcResponse.Reference,
		Quiz: &model.AuthorBasedQuiz{
			QuizItem:      &grpcResponse.Quiz.QuizItem,
			NumberOfItems: &grpcResponse.Quiz.NumberOfItems,
		},
	}

	for _, option := range grpcResponse.Quiz.Options {
		quizResponse.Quiz.Options = append(quizResponse.Quiz.Options, &model.AuthorBasedOptions{QuizWord: &option.QuizWord})
	}

	for _, opt := range grpcResponse.GrammarQuiz {
		grammarQuiz := &model.GrammarQuizAdded{
			CorrectAnswer:    &opt.CorrectAnswer,
			WordInText:       &opt.WordInText,
			ExtraInformation: &opt.ExtraInformation,
		}

		for _, option := range opt.Options {
			grammarQuiz.Options = append(grammarQuiz.Options, &model.AuthorBasedOptions{QuizWord: &option.QuizWord})
		}

		quizResponse.GrammarQuiz = append(quizResponse.GrammarQuiz, grammarQuiz)
	}

	return quizResponse, nil
}

func (s *SokratesHandler) CheckAuthorBased(request *pbxenofon.AnswerRequest, requestID, sessionId string) (*model.AuthorBasedAnswerResponse, error) {
	authorBasedCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	grpcResponse, err := s.AuthorBasedClient.Answer(authorBasedCtx, request)
	if err != nil {
		return nil, err
	}

	return &model.AuthorBasedAnswerResponse{
		Correct:     &grpcResponse.Correct,
		QuizWord:    &grpcResponse.QuizWord,
		WordsInText: convertStringSliceToPointer(grpcResponse.WordsInText),
	}, nil

}

func convertStringSliceToPointer(strings []string) []*string {
	var ptrSlice []*string
	for _, s := range strings {
		ptrSlice = append(ptrSlice, &s)
	}
	return ptrSlice
}
