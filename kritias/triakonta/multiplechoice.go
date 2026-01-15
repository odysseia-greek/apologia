package triakonta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	koinosv1 "github.com/odysseia-greek/apologia/diotima/gen/go/koinos/v1"
	v1 "github.com/odysseia-greek/apologia/kritias/gen/go/v1"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	THEME            string = "theme"
	SET              string = "set"
	GREENGORDER      string = "gre-eng"
	ENGGREORDER      string = "eng-gre"
	OPTIONSEGMENTKEY string = "archytassavedoptions"
)

func (m *MultipleChoiceServiceImpl) Health(context.Context, *koinosv1.HealthRequest) (*koinosv1.HealthResponse, error) {
	elasticHealth := m.Elastic.Health().Info()
	dbHealth := &koinosv1.DatabaseHealth{
		Healthy:       elasticHealth.Healthy,
		ClusterName:   elasticHealth.ClusterName,
		ServerName:    elasticHealth.ServerName,
		ServerVersion: elasticHealth.ServerVersion,
	}

	return &koinosv1.HealthResponse{
		Healthy:        true,
		Time:           time.Now().String(),
		DatabaseHealth: dbHealth,
		Version:        m.Version,
	}, nil
}

func (m *MultipleChoiceServiceImpl) Options(ctx context.Context, request *koinosv1.OptionsRequest) (*v1.AggregatedOptions, error) {
	var unparsedResponse []byte
	cacheItem, _ := m.Archytas.Read(OPTIONSEGMENTKEY)
	if cacheItem != nil {
		unparsedResponse = cacheItem
	} else {
		query := quizAggregationQuery()

		elasticResponse, err := m.Elastic.Query().MatchRaw(m.Index, query)
		if err != nil {
			return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
		}

		unparsedResponse = elasticResponse
		err = m.Archytas.Set(OPTIONSEGMENTKEY, string(elasticResponse))
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

func (m *MultipleChoiceServiceImpl) Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error) {
	if request.Order == "" {
		request.Order = GREENGORDER
	}

	if request.Order != GREENGORDER && request.Order != ENGGREORDER {
		request.Order = GREENGORDER
	}

	var sessionId string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		headerValue := md.Get(config.SessionIdKey)
		if len(headerValue) > 0 {
			sessionId = headerValue[0]
		}
	}

	segmentKey := fmt.Sprintf("%s+%s", request.Theme, request.Set)
	if request.ResetProgress {
		m.Progress.ClearSegment(sessionId, segmentKey)
	}

	if request.ArchiveProgress {
		m.Progress.ResetSegment(sessionId, segmentKey)
	}

	cacheItem, _ := m.Archytas.Read(segmentKey)

	var option models.MediaQuiz

	if cacheItem != nil {
		err := json.Unmarshal(cacheItem, &option)
		if err != nil {
			return nil, err
		}

		go comedy.CacheSpan(string(cacheItem), segmentKey, ctx, m.Streamer)
	} else {
		mustQuery := []map[string]string{
			{
				THEME: request.Theme,
			},
			{
				SET: request.Set,
			},
		}

		query := m.Elastic.Builder().MultipleMatch(mustQuery)
		elasticResponse, err := m.Elastic.Query().Match(m.Index, query)
		if err != nil {
			return nil, err
		}

		if elasticResponse.Hits.Hits == nil || len(elasticResponse.Hits.Hits) == 0 {
			return nil, errors.New("no hits found in query")
		}

		go comedy.DatabaseSpan(query, elasticResponse.Hits.Total.Value, elasticResponse.Took, ctx, m.Streamer)
		source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
		err = json.Unmarshal(source, &option)
		if err != nil {
			return nil, err
		}

		err = m.Archytas.Set(segmentKey, string(source))
		if err != nil {
			if err.Error() != "Key not found" {
				logging.Error(fmt.Sprintf("error when writing cache: %s", err.Error()))
			} else {
				logging.Debug(fmt.Sprintf("cache hit: %s", segmentKey))
			}
		}
	}

	// Ensure session progress is initialized for this segment
	if !m.Progress.Exists(sessionId, segmentKey) || request.ResetProgress || request.ArchiveProgress {
		allGreekWords := make([]string, len(option.Content))
		for i, c := range option.Content {
			allGreekWords[i] = c.Greek
		}
		m.Progress.InitWordsForSegment(sessionId, segmentKey, allGreekWords)
	}

	quiz := &v1.QuizResponse{
		NumberOfItems: int32(len(option.Content)),
	}

	unplayed, unmastered := m.Progress.GetPlayableWords(sessionId, segmentKey, int(request.DoneAfter))
	var wordPool map[string]struct{}

	switch {
	case len(unplayed) > 0:
		wordPool = sliceToSet(unplayed)
	case len(unmastered) > 0:
		wordPool = sliceToSet(unmastered)
	default:
		retryable := m.Progress.GetRetryableWords(sessionId, segmentKey, int(request.DoneAfter))
		wordPool = sliceToSet(retryable)
	}

	var filteredContent []models.MediaContent
	for _, content := range option.Content {
		if _, ok := wordPool[content.Greek]; ok {
			filteredContent = append(filteredContent, content)
		}
	}

	if len(filteredContent) == 0 {
		return nil, status.Errorf(codes.NotFound, "no content available after progress reset")
	}

	var translation string
	var question models.MediaContent
	if len(filteredContent) == 1 {
		question = filteredContent[0]
	} else {
		randNumber := m.Randomizer.RandomNumberBaseZero(len(filteredContent))
		question = filteredContent[randNumber]
	}

	quiz.QuizItem = question.Greek
	translation = question.Translation
	quiz.Options = append(quiz.Options, &v1.Options{
		Option: question.Translation,
	})

	numberOfNeededAnswers := 4

	if len(option.Content) < numberOfNeededAnswers {
		numberOfNeededAnswers = len(option.Content)
	}

	for len(quiz.Options) != numberOfNeededAnswers {
		randNumber := m.Randomizer.RandomNumberBaseZero(len(option.Content))
		randEntry := option.Content[randNumber]

		exists := findQuizWord(quiz.Options, randEntry.Translation)
		if !exists {
			option := &v1.Options{
				Option: randEntry.Translation,
			}
			quiz.Options = append(quiz.Options, option)
		}
	}

	m.Progress.RecordWordPlay(sessionId, segmentKey, quiz.QuizItem, translation)

	rand.Shuffle(len(quiz.Options), func(i, j int) {
		quiz.Options[i], quiz.Options[j] = quiz.Options[j], quiz.Options[i]
	})

	if sessionId != "" {
		progressList, _ := m.Progress.GetProgressForSegment(sessionId, segmentKey, int(request.DoneAfter))
		for word, p := range progressList {
			quiz.Progress = append(quiz.Progress, &koinosv1.ProgressEntry{
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

func (m *MultipleChoiceServiceImpl) Answer(ctx context.Context, request *v1.AnswerRequest) (*v1.AnswerResponse, error) {
	var sessionId string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		headerValue := md.Get(config.SessionIdKey)
		if len(headerValue) > 0 {
			sessionId = headerValue[0]
		}
	}
	segmentKey := fmt.Sprintf("%s+%s", request.Theme, request.Set)
	cacheItem, _ := m.Archytas.Read(segmentKey)

	var option models.MediaQuiz

	if cacheItem != nil {
		err := json.Unmarshal(cacheItem, &option)
		if err != nil {
			return nil, err
		}

		go comedy.CacheSpan(string(cacheItem), segmentKey, ctx, m.Streamer)
	} else {
		mustQuery := []map[string]string{
			{
				THEME: request.Theme,
			},
			{
				SET: request.Set,
			},
		}

		query := m.Elastic.Builder().MultipleMatch(mustQuery)
		elasticResponse, err := m.Elastic.Query().Match(m.Index, query)
		if err != nil {
			return nil, err
		}
		if len(elasticResponse.Hits.Hits) == 0 {
			return nil, fmt.Errorf("no hits found in Elastic")
		}

		go comedy.DatabaseSpan(query, elasticResponse.Hits.Total.Value, elasticResponse.Took, ctx, m.Streamer)

		source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
		err = json.Unmarshal(source, &option)
		if err != nil {
			return nil, err
		}
	}

	answer := v1.AnswerResponse{Correct: false, QuizWord: request.QuizWord}

	for _, content := range option.Content {
		if content.Greek == request.QuizWord {
			if content.Translation == request.Answer {
				answer.Correct = true
			}
			break
		}
	}

	m.Progress.RecordAnswerResult(sessionId, segmentKey, request.QuizWord, answer.Correct)

	if sessionId != "" {
		progressList, finished := m.Progress.GetProgressForSegment(sessionId, segmentKey, int(request.DoneAfter))
		answer.Finished = finished
		for word, p := range progressList {
			answer.Progress = append(answer.Progress, &koinosv1.ProgressEntry{
				Greek:          word,
				Translation:    p.Translation,
				PlayCount:      int32(p.PlayCount),
				CorrectCount:   int32(p.CorrectCount),
				IncorrectCount: int32(p.IncorrectCount),
				LastPlayed:     p.LastPlayed.Format(time.RFC3339),
			})
		}

		if finished {
			m.Progress.ClearSegment(sessionId, segmentKey)
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
func findQuizWord(slice []*v1.Options, val string) bool {
	for _, item := range slice {
		if item.Option == val {
			return true
		}
	}
	return false
}
