package anabasis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	pb "github.com/odysseia-greek/apologia/xenofon/proto"
	"google.golang.org/grpc/metadata"
	"math/rand/v2"
	"time"
)

const (
	THEME            string = "theme"
	SET              string = "set"
	SEGMENT          string = "segment"
	OPTIONSEGMENTKEY string = "archytassavedoptions"
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
		Healthy:        true,
		Time:           time.Now().String(),
		DatabaseHealth: dbHealth,
		Version:        a.Version,
	}, nil
}

func (a *AuthorBasedServiceImpl) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
	var unparsedResponse []byte
	cacheItem, _ := a.Archytas.Read(OPTIONSEGMENTKEY)
	if cacheItem != nil {
		unparsedResponse = cacheItem
	} else {
		query := quizAggregationQuery()

		elasticResponse, err := a.Elastic.Query().MatchRaw(a.Index, query)
		if err != nil {
			return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
		}

		unparsedResponse = elasticResponse
		err = a.Archytas.Set(OPTIONSEGMENTKEY, string(elasticResponse))
		if err != nil {
			logging.Error(err.Error())
		}
	}

	result, err := parseAggregationResult(unparsedResponse)
	if err != nil {
		return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
	}

	return result, nil
}

func (a *AuthorBasedServiceImpl) Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error) {
	var sessionId string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		headerValue := md.Get(config.SessionIdKey)
		if len(headerValue) > 0 {
			sessionId = headerValue[0]
		}
	}

	segmentKey := fmt.Sprintf("%s+%s+%s", request.Theme, request.Set, request.Segment)
	if request.ResetProgress {
		a.Progress.ClearSegment(sessionId, segmentKey)
	}

	if request.ArchiveProgress {
		a.Progress.ResetSegment(sessionId, segmentKey)
	}

	cacheItem, _ := a.Archytas.Read(segmentKey)

	var option models.AuthorbasedQuiz

	if cacheItem != nil {
		err := json.Unmarshal(cacheItem, &option)
		if err != nil {
			return nil, err
		}

		go cacheSpan(string(cacheItem), segmentKey, ctx)
	} else {
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

		source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
		err = json.Unmarshal(source, &option)
		if err != nil {
			return nil, err
		}

		go databaseSpan(elasticResponse, query, ctx)

		err = a.Archytas.Set(segmentKey, string(source))
		if err != nil {
			if err.Error() != "Key not found" {
				logging.Error(fmt.Sprintf("error when writing cache: %s", err.Error()))
			} else {
				logging.Debug(fmt.Sprintf("cache hit: %s", segmentKey))
			}
		}
	}

	if !a.Progress.Exists(sessionId, segmentKey) {
		allGreekWords := make([]string, len(option.Content))
		for i, c := range option.Content {
			allGreekWords[i] = c.Greek
		}
		a.Progress.InitWordsForSegment(sessionId, segmentKey, allGreekWords)
	}

	quiz := &pb.Quiz{
		NumberOfItems: int32(len(option.Content)),
	}

	unplayed, unmastered := a.Progress.GetPlayableWords(sessionId, segmentKey, int(request.DoneAfter))
	var wordPool map[string]struct{}

	switch {
	case len(unplayed) > 0:
		wordPool = sliceToSet(unplayed)
	case len(unmastered) > 0:
		wordPool = sliceToSet(unmastered)
	default:
		retryable := a.Progress.GetRetryableWords(sessionId, segmentKey, int(request.DoneAfter))
		wordPool = sliceToSet(retryable)
	}

	var grammarQuiz []*pb.GrammarQuizAdded
	var filteredContent []models.AuthorBasedContent

	for _, content := range option.Content {
		if _, ok := wordPool[content.Greek]; ok {
			filteredContent = append(filteredContent, content)
		}
	}

	var translation string
	if len(filteredContent) == 1 {
		question := filteredContent[0]
		quiz.QuizItem = question.Greek
		translation = question.Translation
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
		translation = question.Translation
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

	a.Progress.RecordWordPlay(sessionId, segmentKey, quiz.QuizItem, translation)

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

	if sessionId != "" {
		progressList, _ := a.Progress.GetProgressForSegment(sessionId, segmentKey, int(request.DoneAfter))
		for word, p := range progressList {
			authorQuiz.Progress = append(authorQuiz.Progress, &pb.ProgressEntry{
				Greek:          word,
				Translation:    p.Translation,
				PlayCount:      int32(p.PlayCount),
				CorrectCount:   int32(p.CorrectCount),
				IncorrectCount: int32(p.IncorrectCount),
				LastPlayed:     p.LastPlayed.Format(time.RFC3339),
			})
		}
	}

	return &authorQuiz, nil
}

func (a *AuthorBasedServiceImpl) Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.AnswerResponse, error) {
	var sessionId string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		headerValue := md.Get(config.SessionIdKey)
		if len(headerValue) > 0 {
			sessionId = headerValue[0]
		}
	}
	segmentKey := fmt.Sprintf("%s+%s+%s", request.Theme, request.Set, request.Segment)
	cacheItem, _ := a.Archytas.Read(segmentKey)

	var option models.AuthorbasedQuiz

	if cacheItem != nil {
		err := json.Unmarshal(cacheItem, &option)
		if err != nil {
			return nil, err
		}

		go cacheSpan(string(cacheItem), segmentKey, ctx)
	} else {
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
			return nil, fmt.Errorf("no hits found in Elastic")
		}

		go databaseSpan(elasticResponse, query, ctx)

		source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
		err = json.Unmarshal(source, &option)
		if err != nil {
			return nil, err
		}
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

	a.Progress.RecordAnswerResult(sessionId, segmentKey, request.QuizWord, answer.Correct)

	if sessionId != "" {
		progressList, finished := a.Progress.GetProgressForSegment(sessionId, segmentKey, int(request.DoneAfter))
		answer.Finished = finished
		for word, p := range progressList {
			answer.Progress = append(answer.Progress, &pb.ProgressEntry{
				Greek:          word,
				Translation:    p.Translation,
				PlayCount:      int32(p.PlayCount),
				CorrectCount:   int32(p.CorrectCount),
				IncorrectCount: int32(p.IncorrectCount),
				LastPlayed:     p.LastPlayed.Format(time.RFC3339),
			})
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

func sliceToSet(words []string) map[string]struct{} {
	set := make(map[string]struct{}, len(words))
	for _, w := range words {
		set[w] = struct{}{}
	}
	return set
}
