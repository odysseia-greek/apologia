package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.70

import (
	"context"
	"github.com/odysseia-greek/agora/plato/config"
	palkibiades "github.com/odysseia-greek/apologia/alkibiades/proto"
	pbantisthenes "github.com/odysseia-greek/apologia/antisthenes/proto"
	pbartrippos "github.com/odysseia-greek/apologia/aristippos/proto"
	pbkritias "github.com/odysseia-greek/apologia/kritias/proto"
	pbkriton "github.com/odysseia-greek/apologia/kriton/proto"
	"github.com/odysseia-greek/apologia/sokrates/graph/model"
	pbxenofon "github.com/odysseia-greek/apologia/xenofon/proto"
)

// Health is the resolver for the health field.
func (r *queryResolver) Health(ctx context.Context) (*model.AggregatedHealthResponse, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)
	return r.Handler.Health(requestID, sessionId)
}

// MediaOptions is the resolver for the mediaOptions field.
func (r *queryResolver) MediaOptions(ctx context.Context) (*model.AggregatedOptions, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)
	return r.Handler.MediaOptions(requestID, sessionId)
}

// MultipleChoiceOptions is the resolver for the multipleChoiceOptions field.
func (r *queryResolver) MultipleChoiceOptions(ctx context.Context) (*model.ThemedOptions, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)
	return r.Handler.MultipleChoiceOptions(requestID, sessionId)
}

// AuthorBasedOptions is the resolver for the authorBasedOptions field.
func (r *queryResolver) AuthorBasedOptions(ctx context.Context) (*model.AggregatedOptions, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)
	return r.Handler.AuthorBasedOptions(requestID, sessionId)
}

// DialogueOptions is the resolver for the dialogueOptions field.
func (r *queryResolver) DialogueOptions(ctx context.Context) (*model.ThemedOptions, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)
	return r.Handler.DialogueOptions(requestID, sessionId)
}

// GrammarOptions is the resolver for the grammarOptions field.
func (r *queryResolver) GrammarOptions(ctx context.Context) (*model.GrammarOptions, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)
	return r.Handler.GrammarOptions(requestID, sessionId)
}

// JourneyOptions is the resolver for the journeyOptions field.
func (r *queryResolver) JourneyOptions(ctx context.Context) (*model.JourneyOptions, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)
	return r.Handler.JourneyOptions(requestID, sessionId)
}

// MediaAnswer is the resolver for the mediaAnswer field.
func (r *queryResolver) MediaAnswer(ctx context.Context, input *model.MediaAnswerInput) (*model.ComprehensiveResponse, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)

	pb := &pbartrippos.AnswerRequest{
		Theme:         *input.Theme,
		Set:           *input.Set,
		Segment:       *input.Segment,
		Comprehensive: *input.Comprehensive,
		Answer:        *input.Answer,
		QuizWord:      *input.QuizWord,
		DoneAfter:     *input.DoneAfter,
	}

	return r.Handler.CheckMedia(pb, requestID, sessionId)
}

// MediaQuiz is the resolver for the mediaQuiz field.
func (r *queryResolver) MediaQuiz(ctx context.Context, input *model.MediaQuizInput) (*model.MediaQuizResponse, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)

	pb := &pbartrippos.CreationRequest{
		Theme:           *input.Theme,
		Set:             *input.Set,
		Segment:         *input.Segment,
		DoneAfter:       *input.DoneAfter,
		ResetProgress:   *input.ResetProgress,
		ArchiveProgress: *input.ArchiveProgress,
	}

	if input.Order != nil {
		pb.Order = *input.Order
	}

	return r.Handler.CreateMediaQuiz(pb, requestID, sessionId)
}

// MultipleChoiceAnswer is the resolver for the multipleChoiceAnswer field.
func (r *queryResolver) MultipleChoiceAnswer(ctx context.Context, input *model.MultipleChoiceAnswerInput) (*model.ComprehensiveResponse, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)

	pb := &pbkritias.AnswerRequest{
		Theme:         *input.Theme,
		Set:           *input.Set,
		Comprehensive: *input.Comprehensive,
		Answer:        *input.Answer,
		QuizWord:      *input.QuizWord,
	}

	return r.Handler.CheckMultipleChoice(pb, requestID, sessionId)
}

// MultipleChoiceQuiz is the resolver for the multipleChoiceQuiz field.
func (r *queryResolver) MultipleChoiceQuiz(ctx context.Context, input *model.MultipleQuizInput) (*model.MultipleChoiceResponse, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)

	pb := &pbkritias.CreationRequest{
		Theme:           *input.Theme,
		Set:             *input.Set,
		DoneAfter:       *input.DoneAfter,
		ResetProgress:   *input.ResetProgress,
		ArchiveProgress: *input.ArchiveProgress,
	}

	if input.Order != nil {
		pb.Order = *input.Order
	}

	return r.Handler.CreateMultipleChoiceQuiz(pb, requestID, sessionId)
}

// AuthorBasedAnswer is the resolver for the authorBasedAnswer field.
func (r *queryResolver) AuthorBasedAnswer(ctx context.Context, input *model.AuthorBasedAnswerInput) (*model.AuthorBasedAnswerResponse, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)

	pb := &pbxenofon.AnswerRequest{
		Theme:     *input.Theme,
		Set:       *input.Set,
		Segment:   *input.Segment,
		Answer:    *input.Answer,
		QuizWord:  *input.QuizWord,
		DoneAfter: *input.DoneAfter,
	}

	return r.Handler.CheckAuthorBased(pb, requestID, sessionId)
}

// AuthorBasedQuiz is the resolver for the authorBasedQuiz field.
func (r *queryResolver) AuthorBasedQuiz(ctx context.Context, input *model.AuthorBasedInput) (*model.AuthorBasedResponse, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)

	pb := &pbxenofon.CreationRequest{
		Theme:           *input.Theme,
		Set:             *input.Set,
		Segment:         *input.Segment,
		DoneAfter:       *input.DoneAfter,
		ResetProgress:   *input.ResetProgress,
		ArchiveProgress: *input.ArchiveProgress,
	}

	return r.Handler.CreateAuthorBasedQuiz(pb, requestID, sessionId)
}

// AuthorBasedWordForms is the resolver for the authorBasedWordForms field.
func (r *queryResolver) AuthorBasedWordForms(ctx context.Context, input *model.AuthorBasedWordFormsInput) (*model.AuthorBasedWordFormsResponse, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)

	pb := &pbxenofon.WordFormRequest{
		Theme:   *input.Theme,
		Set:     *input.Set,
		Segment: *input.Segment,
	}

	return r.Handler.AuthorBasedWordForms(pb, requestID, sessionId)
}

// DialogueAnswer is the resolver for the dialogueAnswer field.
func (r *queryResolver) DialogueAnswer(ctx context.Context, input *model.DialogueAnswerInput) (*model.DialogueAnswer, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)

	pb := &pbkriton.AnswerRequest{
		Theme: *input.Theme,
		Set:   *input.Set,
	}

	for _, content := range input.Content {
		pb.Content = append(pb.Content, &pbkriton.DialogueContent{
			Translation: *content.Translation,
			Greek:       *content.Greek,
			Place:       *content.Place,
			Speaker:     *content.Speaker,
		})
	}

	return r.Handler.CheckDialogueQuiz(pb, requestID, sessionId)
}

// DialogueQuiz is the resolver for the dialogueQuiz field.
func (r *queryResolver) DialogueQuiz(ctx context.Context, input *model.DialogueQuizInput) (*model.DialogueQuizResponse, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)

	pb := &pbkriton.CreationRequest{
		Theme: *input.Theme,
		Set:   *input.Set,
	}

	return r.Handler.CreateDialogueQuiz(pb, requestID, sessionId)
}

// GrammarQuiz is the resolver for the grammarQuiz field.
func (r *queryResolver) GrammarQuiz(ctx context.Context, input *model.GrammarQuizInput) (*model.GrammarQuizResponse, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)

	pb := &pbantisthenes.CreationRequest{
		Theme:           *input.Theme,
		Set:             *input.Set,
		Segment:         *input.Segment,
		DoneAfter:       *input.DoneAfter,
		ResetProgress:   *input.ResetProgress,
		ArchiveProgress: *input.ArchiveProgress,
	}

	return r.Handler.CreateGrammarQuiz(pb, requestID, sessionId)
}

// GrammarAnswer is the resolver for the grammarAnswer field.
func (r *queryResolver) GrammarAnswer(ctx context.Context, input *model.GrammarAnswerInput) (*model.GrammarAnswer, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)

	pb := &pbantisthenes.AnswerRequest{
		Theme:          *input.Theme,
		Set:            *input.Set,
		Segment:        *input.Segment,
		Answer:         *input.Answer,
		QuizWord:       *input.QuizWord,
		DoneAfter:      *input.DoneAfter,
		DictionaryForm: *input.DictionaryForm,
		Comprehensive:  *input.Comprehensive,
	}

	return r.Handler.CheckGrammar(pb, requestID, sessionId)
}

// JourneyQuiz is the resolver for the journeyQuiz field.
func (r *queryResolver) JourneyQuiz(ctx context.Context, input *model.JourneyQuizInput) (*model.JourneySegmentQuiz, error) {
	requestID, _ := ctx.Value(config.HeaderKey).(string)
	sessionId, _ := ctx.Value(config.SessionIdKey).(string)

	pb := &palkibiades.CreationRequest{
		Theme:   *input.Theme,
		Segment: *input.Segment,
	}

	return r.Handler.CreateJourneySection(requestID, sessionId, pb)
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
