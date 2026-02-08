package gateway

import (
	"context"

	v1 "github.com/odysseia-greek/apologia/alkibiades/gen/go/v1"
	"github.com/odysseia-greek/apologia/alkibiades/strategos"
	koinosv1 "github.com/odysseia-greek/apologia/diotima/gen/go/koinos/v1"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
)

func (s *SokratesHandler) JourneyOptions(ctx context.Context) (*model.JourneyOptions, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.AggregatedOptions

	err := s.JourneyClient.CallWithReconnect(func(client *strategos.JourneyClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Options(outCtx, &koinosv1.OptionsRequest{})
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

func (s *SokratesHandler) CreateJourneySection(ctx context.Context, request *v1.CreationRequest) (*model.JourneySegmentQuiz, error) {
	outCtx, cancel := s.outgoingCtx(ctx)
	defer cancel()

	var grpcResponse *v1.QuizResponse

	err := s.JourneyClient.CallWithReconnect(func(client *strategos.JourneyClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Question(outCtx, request)
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
		case *v1.QuizStep_Match:
			quiz = append(quiz, &model.MatchQuiz{
				Instruction: s.Match.Instruction,
				Pairs:       mapPairs(s.Match.Pairs),
			})
		case *v1.QuizStep_Trivia:
			quiz = append(quiz, &model.TriviaQuiz{
				Question: s.Trivia.Question,
				Options:  s.Trivia.Options,
				Answer:   s.Trivia.Answer,
				Note:     &s.Trivia.Note,
			})
		case *v1.QuizStep_Structure:
			quiz = append(quiz, &model.StructureQuiz{
				Title:    s.Structure.Title,
				Text:     s.Structure.Text,
				Question: s.Structure.Question,
				Options:  s.Structure.Options,
				Answer:   s.Structure.Answer,
				Note:     &s.Structure.Note,
			})
		case *v1.QuizStep_Media:
			quiz = append(quiz, &model.MediaQuiz{
				Instruction: s.Media.Instruction,
				MediaFiles:  mapMediaPairs(s.Media.MediaFiles),
			})
		case *v1.QuizStep_FinalTranslation:
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

func mapPairs(grpcPairs []*v1.MatchPair) []*model.QuizPair {
	var pairs []*model.QuizPair
	for _, p := range grpcPairs {
		pairs = append(pairs, &model.QuizPair{
			Greek:  p.Greek,
			Answer: p.Answer,
		})
	}
	return pairs
}

func mapMediaPairs(files []*v1.MediaEntry) []*model.MediaPair {
	var result []*model.MediaPair
	for _, f := range files {
		result = append(result, &model.MediaPair{
			Word:   f.Word,
			Answer: f.Answer,
		})
	}
	return result
}
