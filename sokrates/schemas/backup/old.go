package schemas

//
//import (
//	"fmt"
//	"github.com/graphql-go/graphql"
//	plato "github.com/odysseia-greek/agora/plato/config"
//	"github.com/odysseia-greek/agora/plato/models"
//)
//
//	// Sokrates
//	"quiz": &graphql.Field{
//		Type: graphql.NewUnion(graphql.UnionConfig{
//			Name:  "QuizResponseUnion",
//			Types: []*graphql.Object{quizResponseType, dialogueQuizType, backup.authorBasedQuizType},
//			ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
//				if _, ok := p.Value.(*models.QuizResponse); ok {
//					return backup.quizResponseType
//				}
//				if _, ok := p.Value.(*models.DialogueQuiz); ok {
//					return backup.dialogueQuizType
//				}
//				if _, ok := p.Value.(*models.AuthorbasedQuizResponse); ok {
//					return backup.authorBasedQuizType
//				}
//				return nil
//			},
//		}),
//		Args: graphql.FieldConfigArgument{
//			"theme": &graphql.ArgumentConfig{
//				Type: graphql.String,
//			},
//			"set": &graphql.ArgumentConfig{
//				Type: graphql.String,
//			},
//			"segment": &graphql.ArgumentConfig{
//				Type: graphql.String,
//			},
//			"quizType": &graphql.ArgumentConfig{
//				Type: graphql.String,
//			},
//			"order": &graphql.ArgumentConfig{
//				Type: graphql.String,
//			},
//			"excludeWords": &graphql.ArgumentConfig{
//				Type: graphql.NewList(graphql.String),
//			},
//		},
//		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
//			ctx := p.Context
//			traceID, ok := ctx.Value(plato.HeaderKey).(string)
//			if !ok {
//				return nil, errors.New("failed to get request from context")
//			}
//
//			theme, _ := p.Args["theme"].(string)
//			segment, _ := p.Args["segment"].(string)
//			order, _ := p.Args["order"].(string)
//			excludeWords, _ := p.Args["excludeWords"].([]interface{})
//			excludeWordsStr := make([]string, len(excludeWords))
//
//			if excludeWords != nil {
//				for i, word := range excludeWords {
//					excludeWordsStr[i], _ = word.(string)
//				}
//			}
//
//			set, isOK := p.Args["set"].(string)
//			if !isOK {
//				return nil, fmt.Errorf("expected argument set")
//			}
//			quizType, isOK := p.Args["quizType"].(string)
//			if !isOK {
//				return nil, fmt.Errorf("expected argument quizType")
//			}
//
//			if quizType == models.DIALOGUE {
//				return handler.CreateDialogueQuiz(theme, set, segment, quizType, traceID)
//			} else if quizType == models.AUTHORBASED {
//				return handler.CreateAuthorBasedQuiz(theme, set, segment, quizType, traceID, excludeWordsStr)
//			} else {
//				return handler.CreateQuiz(theme, set, segment, quizType, order, traceID, excludeWordsStr)
//			}
//
//		},
//	},
//
//		"answer": &graphql.Field{
//		Type: graphql.NewUnion(graphql.UnionConfig{
//			Name:  "AnswerUnion",
//			Types: []*graphql.Object{backup.comprehensiveAnswer, backup.dialogueAnswerType, backup.authorBasedAnswer},
//			ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
//				if _, ok := p.Value.(*models.ComprehensiveResponse); ok {
//					return backup.comprehensiveAnswer
//				}
//				if _, ok := p.Value.(*models.DialogueAnswer); ok {
//					return backup.dialogueAnswerType
//				}
//				if _, ok := p.Value.(*models.AuthorBasedResponse); ok {
//					return backup.authorBasedAnswer
//				}
//				return nil
//			},
//		}),
//		Args: graphql.FieldConfigArgument{
//			"theme": &graphql.ArgumentConfig{
//				Type: graphql.String,
//			},
//			"set": &graphql.ArgumentConfig{
//				Type: graphql.String,
//			},
//			"segment": &graphql.ArgumentConfig{
//				Type: graphql.String,
//			},
//			"quizType": &graphql.ArgumentConfig{
//				Type: graphql.String,
//			},
//			"quizWord": &graphql.ArgumentConfig{
//				Type: graphql.String,
//			},
//			"answer": &graphql.ArgumentConfig{
//				Type: graphql.String,
//			},
//			"comprehensive": &graphql.ArgumentConfig{
//				Type: graphql.Boolean,
//			},
//			"dialogue": &graphql.ArgumentConfig{
//				Type: graphql.NewList(backup.dialogueInputType),
//			},
//		},
//		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
//			ctx := p.Context
//
//			// Get the traceID
//			traceID, ok := ctx.Value(plato.HeaderKey).(string)
//			if !ok {
//				return nil, errors.New("failed to get request from context")
//			}
//
//			set, isOK := p.Args["set"].(string)
//			if !isOK {
//				return nil, fmt.Errorf("expected argument set")
//			}
//			quizType, isOK := p.Args["quizType"].(string)
//			if !isOK {
//				return nil, fmt.Errorf("expected argument quizType")
//			}
//
//			theme, _ := p.Args["theme"].(string)
//			segment, _ := p.Args["segment"].(string)
//			quizWord, _ := p.Args["quizWord"].(string)
//			answer, _ := p.Args["answer"].(string)
//			comprehensive, _ := p.Args["comprehensive"].(bool)
//			dialogueList, _ := p.Args["dialogue"].([]interface{})
//
//			var dialogue []models.DialogueContent
//			for _, item := range dialogueList {
//				itemMap, ok := item.(map[string]interface{})
//				if !ok {
//					return nil, fmt.Errorf("each dialogue item must be a map")
//				}
//
//				var dialogueItem models.DialogueContent
//				if translation, ok := itemMap["translation"].(string); ok {
//					dialogueItem.Translation = translation
//				}
//				if greek, ok := itemMap["greek"].(string); ok {
//					dialogueItem.Greek = greek
//				}
//				if place, ok := itemMap["place"].(int); ok {
//					dialogueItem.Place = place
//				}
//				if speaker, ok := itemMap["speaker"].(string); ok {
//					dialogueItem.Speaker = speaker
//				}
//
//				dialogue = append(dialogue, dialogueItem)
//			}
//
//			answerRequest := models.AnswerRequest{
//				Theme:         theme,
//				Set:           set,
//				Segment:       segment,
//				QuizType:      quizType,
//				Comprehensive: comprehensive,
//				Answer:        answer,
//				Dialogue:      dialogue,
//				QuizWord:      quizWord,
//			}
//
//			if quizType == models.DIALOGUE {
//				return handler.CheckDialogue(answerRequest, traceID)
//			} else if quizType == models.AUTHORBASED {
//				return handler.CheckAuthorBased(answerRequest, traceID)
//			} else {
//				return handler.Check(answerRequest, traceID)
//			}
//		},
//	},
//
//		"options": &graphql.Field{
//		Type:        backup.aggregateResultType,
//		Description: "returns the options for a specific quiztype",
//		Args: graphql.FieldConfigArgument{
//			"quizType": &graphql.ArgumentConfig{
//				Type: graphql.String,
//			},
//		},
//		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
//			ctx := p.Context
//			// Get the traceID
//			traceID, ok := ctx.Value(plato.HeaderKey).(string)
//			if !ok {
//				return nil, errors.New("failed to get request from context")
//			}
//
//			quizType, isOK := p.Args["quizType"].(string)
//			if !isOK {
//				return nil, fmt.Errorf("expected argument quizType")
//			}
//			return handler.Options(quizType, traceID)
//		},
//	},
//}
