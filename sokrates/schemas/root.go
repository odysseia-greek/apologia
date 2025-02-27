package schemas

import (
	"context"
	"errors"
	"github.com/graphql-go/graphql"
	plato "github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/apologia/sokrates/gateway"
	"log"
	"sync"
)

var (
	handler             *gateway.SokratesHandler
	sokratesHandlerOnce sync.Once
)

func SokratesHandler() *gateway.SokratesHandler {
	sokratesHandlerOnce.Do(func() {
		ctx := context.Background()
		homerosHandler, err := gateway.CreateNewConfig(ctx)
		if err != nil {
			log.Print(err)
		}
		handler = homerosHandler
	})
	return handler
}

var SokratesSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		// MEDIA QUIZ
		"mediaQuiz": &graphql.Field{
			Type: quizResponseType,
			Args: quizArgs(),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				request := buildQuizRequest(p)
				return handler.CreateMediaQuiz(request, traceID)
			},
		},

		//// DIALOGUE QUIZ
		//"dialogueQuiz": &graphql.Field{
		//	Type: dialogueQuizType,
		//	Args: quizArgs(),
		//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		//		ctx := p.Context
		//		traceID, ok := ctx.Value(plato.HeaderKey).(string)
		//		if !ok {
		//			return nil, errors.New("failed to get request from context")
		//		}
		//
		//		request := buildQuizRequest(p)
		//		return handler.CreateDialogueQuiz(request, traceID)
		//	},
		//},
		//
		//// MULTIPLE CHOICE QUIZ
		//"multipleChoiceQuiz": &graphql.Field{
		//	Type: multipleChoiceQuizType,
		//	Args: quizArgs(),
		//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		//		ctx := p.Context
		//		traceID, ok := ctx.Value(plato.HeaderKey).(string)
		//		if !ok {
		//			return nil, errors.New("failed to get request from context")
		//		}
		//
		//		request := buildQuizRequest(p)
		//		return handler.CreateMultipleChoiceQuiz(request, traceID)
		//	},
		//},
		//
		//// AUTHOR-BASED QUIZ
		//"authorBasedQuiz": &graphql.Field{
		//	Type: authorBasedQuizType,
		//	Args: quizArgs(),
		//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		//		ctx := p.Context
		//		traceID, ok := ctx.Value(plato.HeaderKey).(string)
		//		if !ok {
		//			return nil, errors.New("failed to get request from context")
		//		}
		//
		//		request := buildQuizRequest(p)
		//		return handler.CreateAuthorBasedQuiz(request, traceID)
		//	},
		//},
		//
		//// FUTURE QUIZZES
		//"grammarQuiz": &graphql.Field{
		//	Type: grammarQuizType,
		//	Args: quizArgs(),
		//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		//		// Future implementation
		//		return nil, errors.New("not yet implemented")
		//	},
		//},
		//
		//"journeyQuiz": &graphql.Field{
		//	Type: journeyQuizType,
		//	Args: quizArgs(),
		//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		//		// Future implementation
		//		return nil, errors.New("not yet implemented")
		//	},
		//},

		// MEDIA ANSWER
		"mediaAnswer": &graphql.Field{
			Type: comprehensiveAnswer,
			Args: answerArgs(),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context
				traceID, ok := ctx.Value(plato.HeaderKey).(string)
				if !ok {
					return nil, errors.New("failed to get request from context")
				}

				request := buildAnswerRequest(p)
				return handler.CheckMedia(request, traceID)
			},
		},

		//// DIALOGUE ANSWER
		//"dialogueAnswer": &graphql.Field{
		//	Type: dialogueAnswerType,
		//	Args: answerArgs(),
		//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		//		ctx := p.Context
		//		traceID, ok := ctx.Value(plato.HeaderKey).(string)
		//		if !ok {
		//			return nil, errors.New("failed to get request from context")
		//		}
		//
		//		request := buildAnswerRequest(p)
		//		return handler.CheckDialogue(request, traceID)
		//	},
		//},

		//// OPTIONS
		//"options": &graphql.Field{
		//	Type: aggregateResultType,
		//	Args: graphql.FieldConfigArgument{
		//		"quizType": &graphql.ArgumentConfig{
		//			Type: graphql.String,
		//		},
		//	},
		//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		//		ctx := p.Context
		//		traceID, ok := ctx.Value(plato.HeaderKey).(string)
		//		if !ok {
		//			return nil, errors.New("failed to get request from context")
		//		}
		//
		//		quizType, isOK := p.Args["quizType"].(string)
		//		if !isOK {
		//			return nil, fmt.Errorf("expected argument quizType")
		//		}
		//
		//		return handler.Options(quizType, traceID)
		//	},
		//},
	},
})
