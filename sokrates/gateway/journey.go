package gateway

import (
	pbalkibiades "github.com/odysseia-greek/apologia/alkibiades/proto"
	"github.com/odysseia-greek/apologia/alkibiades/strategos"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) JourneyOptions(requestID, sessionId string) (*model.JourneyOptions, error) {
	optionsCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	var grpcResponse *pbalkibiades.AggregatedOptions

	err := s.JourneyClient.CallWithReconnect(func(client *strategos.JourneyClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Options(optionsCtx, &pbalkibiades.OptionsRequest{})
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	var themes []*model.JourneyThemes
	for _, grpcTheme := range grpcResponse.Themes {
		var segments []*model.JourneySegment
		for _, grpcSegment := range grpcTheme.Segments {
			segments = append(segments, &model.JourneySegment{
				Name:     &grpcSegment.Name,
				Number:   &grpcSegment.Number,
				Location: &grpcSegment.Location,
				Coordinates: &model.Coordinates{
					X: float32ToFloat64Ptr(grpcSegment.Coordinates.X),
					Y: float32ToFloat64Ptr(grpcSegment.Coordinates.Y),
				},
			})
		}

		themes = append(themes, &model.JourneyThemes{
			Name:     &grpcTheme.Name,
			Segments: segments,
		})
	}

	return &model.JourneyOptions{
		Themes: themes,
	}, nil
}

func (s *SokratesHandler) CreateJourneySection(requestID, sessionId string, request *pbalkibiades.CreationRequest) (*model.JourneySegmentQuiz, error) {
	journeyCreateCtx, cancel := s.createRequestHeader(requestID, sessionId)
	defer cancel()

	var grpcResponse *pbalkibiades.QuizResponse

	err := s.JourneyClient.CallWithReconnect(func(client *strategos.JourneyClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Question(journeyCreateCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	response := &model.JourneySegmentQuiz{
		Theme:       grpcResponse.Theme,
		Segment:     grpcResponse.Segment,
		Number:      grpcResponse.Number,
		Sentence:    grpcResponse.Sentence,
		Translation: grpcResponse.Translation,
		ContextNote: &grpcResponse.ContextNote,
		Intro:       nil,
		Quiz:        nil,
	}

	var intro *model.QuizIntro
	if grpcResponse.Intro != nil {
		intro = &model.QuizIntro{
			Author:     grpcResponse.Intro.Author,
			Work:       grpcResponse.Intro.Work,
			Background: grpcResponse.Intro.Background,
		}
	}

	response.Intro = intro
	var quiz []model.QuizSection

	for _, section := range grpcResponse.Quiz {
		switch s := section.Type.(type) {
		case *pbalkibiades.QuizStep_Match:
			quiz = append(quiz, &model.MatchQuiz{
				Instruction: s.Match.Instruction,
				Pairs:       mapPairs(s.Match.Pairs),
			})
		case *pbalkibiades.QuizStep_Trivia:
			quiz = append(quiz, &model.TriviaQuiz{
				Question: s.Trivia.Question,
				Options:  s.Trivia.Options,
				Answer:   s.Trivia.Answer,
				Note:     &s.Trivia.Note,
			})
		case *pbalkibiades.QuizStep_Structure:
			quiz = append(quiz, &model.StructureQuiz{
				Title:    s.Structure.Title,
				Text:     s.Structure.Text,
				Question: s.Structure.Question,
				Options:  s.Structure.Options,
				Answer:   s.Structure.Answer,
				Note:     &s.Structure.Note,
			})
		case *pbalkibiades.QuizStep_Media:
			quiz = append(quiz, &model.MediaQuiz{
				Instruction: s.Media.Instruction,
				MediaFiles:  mapMediaPairs(s.Media.MediaFiles),
			})
		case *pbalkibiades.QuizStep_FinalTranslation:
			quiz = append(quiz, &model.FinalTranslationQuiz{
				Instruction: s.FinalTranslation.Instruction,
				Options:     s.FinalTranslation.Options,
				Answer:      s.FinalTranslation.Answer,
			})
		default:
			// optionally log/skip unknown types
		}
	}

	response.Quiz = quiz
	return response, nil
}

func float32ToFloat64Ptr(f float32) *float64 {
	val := float64(f)
	return &val
}

func mapPairs(grpcPairs []*pbalkibiades.MatchPair) []*model.QuizPair {
	var pairs []*model.QuizPair
	for _, p := range grpcPairs {
		pairs = append(pairs, &model.QuizPair{
			Greek:  p.Greek,
			Answer: p.Answer,
		})
	}
	return pairs
}

func mapMediaPairs(files []*pbalkibiades.MediaEntry) []*model.MediaPair {
	var result []*model.MediaPair
	for _, f := range files {
		result = append(result, &model.MediaPair{
			Word:   f.Word,
			Answer: f.Answer,
		})
	}
	return result
}
