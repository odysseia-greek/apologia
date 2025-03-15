package anabasis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	pb "github.com/odysseia-greek/apologia/xenofon/proto"
	"math/rand/v2"
	"time"
)

const (
	THEME   string = "theme"
	SET     string = "set"
	SEGMENT string = "segment"
)

func (a *AuthorBasedServiceImpl) Health(context.Context, *pb.HealthRequest) (*pb.HealthResponse, error) {
	elasticHealth := a.Elastic.Health().Info()
	dbHealth := &pb.DatabaseHealth{
		Healthy:       elasticHealth.Healthy,
		ClusterName:   elasticHealth.ClusterName,
		ServerName:    elasticHealth.ServerName,
		ServerVersion: elasticHealth.ServerVersion,
	}

	return &pb.HealthResponse{
		Healthy:        dbHealth.Healthy,
		Time:           time.Now().String(),
		DatabaseHealth: dbHealth,
	}, nil
}

func (a *AuthorBasedServiceImpl) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
	query := quizAggregationQuery()

	elasticResult, err := a.Elastic.Query().MatchRaw(a.Index, query)
	if err != nil {
		return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
	}

	result, err := parseAggregationResult(elasticResult)
	if err != nil {
		return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
	}

	return result, nil
}

func (a *AuthorBasedServiceImpl) Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error) {
	mustQuery := []map[string]string{
		{
			THEME: request.Theme,
		},
		{
			SET: request.Set,
		},
		{
			SEGMENT: request.Segment,
		},
	}
	query := a.Elastic.Builder().MultipleMatch(mustQuery)
	elasticResponse, err := a.Elastic.Query().Match(a.Index, query)
	if err != nil {
		return nil, err
	}

	if elasticResponse.Hits.Hits == nil || len(elasticResponse.Hits.Hits) == 0 {
		return nil, errors.New("no hits found in query")
	}

	var option models.AuthorbasedQuiz

	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &option)
	if err != nil {
		return nil, err
	}

	//if traceCall {
	//	go a.databaseSpan(elasticResponse, query, traceID, spanID)
	//}

	quiz := &pb.Quiz{
		NumberOfItems: int32(len(option.Content)),
	}

	var grammarQuiz []*pb.GrammarQuizAdded
	var filteredContent []models.AuthorBasedContent

	for _, content := range option.Content {
		addWord := true
		for _, word := range request.ExcludeWords {
			if content.Greek == word {
				addWord = false
			}
		}

		if addWord {
			filteredContent = append(filteredContent, content)
		}
	}

	if len(filteredContent) == 1 {
		question := filteredContent[0]
		quiz.QuizItem = question.Greek
		quiz.Options = append(quiz.Options, &pb.Options{
			QuizWord: question.Translation,
		})
	} else {
		randNumber := a.Randomizer.RandomNumberBaseZero(len(filteredContent))
		question := filteredContent[randNumber]
		if question.HasGrammarQuestions {
			//add grammar question
			for _, grammarQuestion := range question.GrammarQuestions {
				grammarQuizOption := &pb.GrammarQuizAdded{
					CorrectAnswer:    grammarQuestion.CorrectAnswer,
					WordInText:       grammarQuestion.WordInText,
					ExtraInformation: grammarQuestion.ExtraInformation,
					Options:          nil,
				}

				var setToQuery []string
				switch grammarQuestion.TypeOfWord {
				case "noun":
					setToQuery = option.GrammarQuestionOptions.Nouns
				case "verb":
					setToQuery = option.GrammarQuestionOptions.Verbs
				case "misc":
					setToQuery = option.GrammarQuestionOptions.Misc
				default:
					setToQuery = option.GrammarQuestionOptions.Nouns
				}

				var grammarOptions []*pb.Options

				grammarOptions = append(grammarOptions, &pb.Options{
					QuizWord: grammarQuizOption.CorrectAnswer,
				})

				numberOfNeededAnswers := 4
				for len(grammarOptions) != numberOfNeededAnswers {
					randNumber := a.Randomizer.RandomNumberBaseZero(len(setToQuery))
					randEntry := setToQuery[randNumber]

					exists := findQuizWord(grammarOptions, randEntry)
					if !exists {
						grammarOption := &pb.Options{
							QuizWord: randEntry,
						}
						grammarOptions = append(grammarOptions, grammarOption)
					}
				}

				grammarQuizOption.Options = grammarOptions
				rand.Shuffle(len(grammarQuizOption.Options), func(i, j int) {
					grammarQuizOption.Options[i], grammarQuizOption.Options[j] = grammarQuizOption.Options[j], grammarQuizOption.Options[i]
				})
				grammarQuiz = append(grammarQuiz, grammarQuizOption)
			}
		}
		quiz.QuizItem = question.Greek
		quiz.Options = append(quiz.Options, &pb.Options{
			QuizWord: question.Translation,
		})
	}

	numberOfNeededAnswers := 4

	if len(option.Content) < numberOfNeededAnswers {
		numberOfNeededAnswers = len(option.Content)
	}

	for len(quiz.Options) != numberOfNeededAnswers {
		randNumber := a.Randomizer.RandomNumberBaseZero(len(option.Content))
		randEntry := option.Content[randNumber]

		exists := findQuizWord(quiz.Options, randEntry.Translation)
		if !exists {
			option := &pb.Options{
				QuizWord: randEntry.Translation,
			}
			quiz.Options = append(quiz.Options, option)
		}
	}

	rand.Shuffle(len(quiz.Options), func(i, j int) {
		quiz.Options[i], quiz.Options[j] = quiz.Options[j], quiz.Options[i]
	})

	authorQuiz := pb.QuizResponse{
		FullSentence: option.FullSentence,
		Translation:  option.Translation,
		Reference:    option.Reference,
		Quiz:         quiz,
		GrammarQuiz:  grammarQuiz,
	}

	return &authorQuiz, nil
}

func (a *AuthorBasedServiceImpl) Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.AnswerResponse, error) {
	mustQuery := []map[string]string{
		{
			THEME: request.Theme,
		},
		{
			SEGMENT: request.Segment,
		},
		{
			SET: request.Set,
		},
	}

	query := a.Elastic.Builder().MultipleMatch(mustQuery)
	elasticResponse, err := a.Elastic.Query().Match(a.Index, query)
	if err != nil {
		return nil, err
	}
	if len(elasticResponse.Hits.Hits) == 0 {
		logging.Error(fmt.Sprintf("no hits found in Elastic"))
		return nil, fmt.Errorf("no hits found in Elastic")
	}

	var option models.AuthorbasedQuiz
	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &option)
	if err != nil {
		return nil, err
	}

	answer := &pb.AnswerResponse{
		Correct:  false,
		QuizWord: request.QuizWord,
	}

	for _, content := range option.Content {
		if content.Greek == request.QuizWord {
			if content.Translation == request.Answer {
				answer.Correct = true
				answer.WordsInText = content.WordsInText
			}
		}
	}

	return answer, nil
}

// findQuizWord takes a slice and looks for an element in it
func findQuizWord(slice []*pb.Options, val string) bool {
	for _, item := range slice {
		if item.QuizWord == val {
			return true
		}
	}
	return false
}
