package kunismos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	v1 "github.com/odysseia-greek/apologia/antisthenes/gen/go/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	OPTIONSEGMENTKEY string = "archytassavedoptions"
	THEME            string = "theme"
	SET              string = "set"
	SEGMENT          string = "segment"
)

func (g *GrammarServiceImpl) Health(context.Context, *v1.HealthRequest) (*v1.HealthResponse, error) {
	elasticHealth := g.Elastic.Health().Info()
	dbHealth := &v1.DatabaseHealth{
		Healthy:       elasticHealth.Healthy,
		ClusterName:   elasticHealth.ClusterName,
		ServerName:    elasticHealth.ServerName,
		ServerVersion: elasticHealth.ServerVersion,
	}

	return &v1.HealthResponse{
		Healthy:        true,
		Time:           time.Now().String(),
		DatabaseHealth: dbHealth,
		Version:        g.Version,
	}, nil
}

func (g *GrammarServiceImpl) Options(ctx context.Context, request *v1.OptionsRequest) (*v1.AggregatedOptions, error) {
	var unparsedResponse []byte
	cacheItem, _ := g.Archytas.Read(OPTIONSEGMENTKEY)
	if cacheItem != nil {
		unparsedResponse = cacheItem
	} else {
		query := quizAggregationQuery()

		elasticResponse, err := g.Elastic.Query().MatchRaw(g.Index, query)
		if err != nil {
			return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
		}

		unparsedResponse = elasticResponse
		err = g.Archytas.Set(OPTIONSEGMENTKEY, string(elasticResponse))
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

func (g *GrammarServiceImpl) Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error) {
	var sessionId string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		headerValue := md.Get(config.SessionIdKey)
		if len(headerValue) > 0 {
			sessionId = headerValue[0]
		}
	}

	rawKey := fmt.Sprintf("%s+%s+%s", request.Theme, request.Set, request.Segment)
	segmentKey := strings.ReplaceAll(rawKey, " ", "")

	if request.ResetProgress {
		g.Progress.ClearSegment(sessionId, segmentKey)
	}

	if request.ArchiveProgress {
		g.Progress.ResetSegment(sessionId, segmentKey)
	}

	cacheItem, _ := g.Archytas.Read(segmentKey)

	var option GrammarBasedQuiz

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

		query := g.Elastic.Builder().MultipleMatch(mustQuery)
		elasticResponse, err := g.Elastic.Query().Match(g.Index, query)
		if err != nil {
			return nil, err
		}

		if elasticResponse.Hits.Hits == nil || len(elasticResponse.Hits.Hits) == 0 {
			return nil, errors.New("no hits found in query")
		}

		go databaseSpan(elasticResponse, query, ctx)
		source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
		err = json.Unmarshal(source, &option)
		if err != nil {
			return nil, err
		}

		err = g.Archytas.Set(segmentKey, string(source))
		if err != nil {
			if err.Error() != "Key not found" {
				logging.Error(fmt.Sprintf("error when writing cache: %s", err.Error()))
			} else {
				logging.Debug(fmt.Sprintf("cache hit: %s", segmentKey))
			}
		}
	}

	// Ensure session progress is initialized for this segment
	if !g.Progress.Exists(sessionId, segmentKey) || request.ResetProgress || request.ArchiveProgress {
		logging.Info(fmt.Sprintf("initializing progress for segment: %s and session: %s", segmentKey, sessionId))
		allGreekWords := make([]string, len(option.Content))
		for i, c := range option.Content {
			allGreekWords[i] = c.Greek
		}
		g.Progress.InitWordsForSegment(sessionId, segmentKey, allGreekWords)
	}

	quiz := &v1.QuizResponse{
		NumberOfItems:   int32(len(option.Content)),
		Description:     option.Description,
		Difficulty:      option.Difficulty,
		ContractionRule: option.ContractionRule,
	}

	unplayed, unmastered := g.Progress.GetPlayableWords(sessionId, segmentKey, int(request.DoneAfter))
	var wordPool map[string]struct{}

	switch {
	case len(unplayed) > 0:
		wordPool = sliceToSet(unplayed)
	case len(unmastered) > 0:
		wordPool = sliceToSet(unmastered)
	default:
		retryable := g.Progress.GetRetryableWords(sessionId, segmentKey, int(request.DoneAfter))
		wordPool = sliceToSet(retryable)
	}

	var filteredContent []GrammarContent
	for _, content := range option.Content {
		if _, ok := wordPool[content.Greek]; ok {
			filteredContent = append(filteredContent, content)
		}
	}

	if len(filteredContent) == 0 {
		return nil, status.Errorf(codes.NotFound, "no content available after progress reset")
	}

	var correctAnswer string
	var question GrammarContent

	if len(filteredContent) == 1 {
		question = filteredContent[0]

	} else {
		randNumber := g.Randomizer.RandomNumberBaseZero(len(filteredContent))
		question = filteredContent[randNumber]
	}

	quiz.QuizItem = question.Greek
	quiz.DictionaryForm = question.DictionaryForm
	quiz.Stem = question.Stem
	quiz.Translation = question.Translation
	quiz.Options = append(quiz.Options, &v1.GrammarOptions{
		Option: question.GrammarQuestion.CorrectAnswer,
	})

	correctAnswer = question.GrammarQuestion.CorrectAnswer

	numberOfNeededAnswers := 4

	if len(option.Content) < numberOfNeededAnswers {
		numberOfNeededAnswers = len(option.Content)
	}

	for len(quiz.Options) != numberOfNeededAnswers {
		randNumber := g.Randomizer.RandomNumberBaseZero(len(option.Content))
		randEntry := option.Content[randNumber]

		exists := findQuizWord(quiz.Options, randEntry.GrammarQuestion.CorrectAnswer)
		if !exists {
			quiz.Options = append(quiz.Options, &v1.GrammarOptions{
				Option: randEntry.GrammarQuestion.CorrectAnswer,
			})
		}
	}

	g.Progress.RecordWordPlay(sessionId, segmentKey, quiz.QuizItem, correctAnswer)

	rand.Shuffle(len(quiz.Options), func(i, j int) {
		quiz.Options[i], quiz.Options[j] = quiz.Options[j], quiz.Options[i]
	})

	if sessionId != "" {
		progressList, _ := g.Progress.GetProgressForSegment(sessionId, segmentKey, int(request.DoneAfter))
		for word, p := range progressList {
			quiz.Progress = append(quiz.Progress, &v1.ProgressEntry{
				Greek:          word,
				Translation:    p.Translation,
				PlayCount:      int32(p.PlayCount),
				CorrectCount:   int32(p.CorrectCount),
				IncorrectCount: int32(p.IncorrectCount),
				LastPlayed:     p.LastPlayed.Format(time.RFC3339),
			})
		}
	}

	return quiz, nil
}

func (g *GrammarServiceImpl) Answer(ctx context.Context, request *v1.AnswerRequest) (*v1.AnswerResponse, error) {
	var sessionId string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		headerValue := md.Get(config.SessionIdKey)
		if len(headerValue) > 0 {
			sessionId = headerValue[0]
		}
	}

	rawKey := fmt.Sprintf("%s+%s+%s", request.Theme, request.Set, request.Segment)
	segmentKey := strings.ReplaceAll(rawKey, " ", "")
	cacheItem, _ := g.Archytas.Read(segmentKey)

	var option GrammarBasedQuiz

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

		query := g.Elastic.Builder().MultipleMatch(mustQuery)
		elasticResponse, err := g.Elastic.Query().Match(g.Index, query)
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

	answer := v1.AnswerResponse{Correct: false, QuizWord: request.DictionaryForm}

	for _, content := range option.Content {
		if content.Greek == request.QuizWord {
			if content.GrammarQuestion.CorrectAnswer == request.Answer {
				answer.Correct = true
			}
			break
		}
	}

	g.Progress.RecordAnswerResult(sessionId, segmentKey, request.QuizWord, answer.Correct)

	if sessionId != "" {
		progressList, finished := g.Progress.GetProgressForSegment(sessionId, segmentKey, int(request.DoneAfter))
		answer.Finished = finished
		for word, p := range progressList {
			answer.Progress = append(answer.Progress, &v1.ProgressEntry{
				Greek:          word,
				Translation:    p.Translation,
				PlayCount:      int32(p.PlayCount),
				CorrectCount:   int32(p.CorrectCount),
				IncorrectCount: int32(p.IncorrectCount),
				LastPlayed:     p.LastPlayed.Format(time.RFC3339),
			})
		}

		if finished {
			g.Progress.ClearSegment(sessionId, segmentKey)

			var greekWords []string
			for _, content := range option.Content {
				greekWords = append(greekWords, content.Greek)
			}
			g.Progress.InitWordsForSegment(sessionId, segmentKey, greekWords)
		}
	}

	return &answer, nil
}

func sliceToSet(words []string) map[string]struct{} {
	set := make(map[string]struct{}, len(words))
	for _, w := range words {
		set[w] = struct{}{}
	}
	return set
}

// findQuizWord takes a slice and looks for an element in it
func findQuizWord(slice []*v1.GrammarOptions, val string) bool {
	for _, item := range slice {
		if item.Option == val {
			return true
		}
	}
	return false
}
