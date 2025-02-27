package hedone

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/agora/plato/transform"
	pb "github.com/odysseia-greek/apologia/aristippos/proto"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc/metadata"
	"math/rand/v2"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	THEME       string = "theme"
	SET         string = "set"
	SEGMENT     string = "segment"
	MEDIA       string = "media"
	QUIZTYPE    string = "quizType"
	GREENGORDER string = "gre-eng"
	ENGGREORDER string = "eng-gre"
)

func (m *MediaServiceImpl) Health(context.Context, *pb.HealthRequest) (*pb.HealthResponse, error) {
	elasticHealth := m.Elastic.Health().Info()
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

func (m *MediaServiceImpl) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
	query := quizAggregationQuery()

	elasticResult, err := m.Elastic.Query().MatchRaw(m.Index, query)
	if err != nil {
		return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
	}

	result, err := parseAggregationResult(elasticResult)
	if err != nil {
		return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
	}

	return result, nil
}

func (m *MediaServiceImpl) Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error) {
	if request.Order == "" {
		request.Order = GREENGORDER
	}

	if request.Order != GREENGORDER && request.Order != ENGGREORDER {
		request.Order = GREENGORDER
	}

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
		{
			QUIZTYPE: MEDIA,
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

	var option models.MediaQuiz

	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &option)
	if err != nil {
		return nil, err
	}

	quiz := &pb.QuizResponse{
		NumberOfItems: int32(len(option.Content)),
	}

	var filteredContent []models.MediaContent

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
			Option:   question.Translation,
			ImageUrl: question.ImageURL,
		})
	} else {
		randNumber := m.Randomizer.RandomNumberBaseZero(len(filteredContent))
		question := filteredContent[randNumber]
		quiz.QuizItem = question.Greek
		quiz.Options = append(quiz.Options, &pb.Options{
			Option:   question.Translation,
			ImageUrl: question.ImageURL,
		})
	}

	numberOfNeededAnswers := 4

	if len(option.Content) < numberOfNeededAnswers {
		numberOfNeededAnswers = len(option.Content)
	}

	for len(quiz.Options) != numberOfNeededAnswers {
		randNumber := m.Randomizer.RandomNumberBaseZero(len(option.Content))
		randEntry := option.Content[randNumber]

		exists := findQuizWord(quiz.Options, randEntry.Translation)
		if !exists {
			option := &pb.Options{
				Option:   randEntry.Translation,
				ImageUrl: randEntry.ImageURL,
			}
			quiz.Options = append(quiz.Options, option)
		}
	}

	rand.Shuffle(len(quiz.Options), func(i, j int) {
		quiz.Options[i], quiz.Options[j] = quiz.Options[j], quiz.Options[i]
	})

	return quiz, nil
}

func (m *MediaServiceImpl) Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.ComprehensiveResponse, error) {
	mustQuery := []map[string]string{
		{
			QUIZTYPE: MEDIA,
		},
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

	query := m.Elastic.Builder().MultipleMatch(mustQuery)
	elasticResponse, err := m.Elastic.Query().Match(m.Index, query)
	if err != nil {
		return nil, err
	}
	if len(elasticResponse.Hits.Hits) == 0 {
		return nil, fmt.Errorf("no hits found in Elastic")
	}

	var option models.MediaQuiz
	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &option)
	if err != nil {
		return nil, err
	}

	answer := pb.ComprehensiveResponse{Correct: false, QuizWord: request.QuizWord}

	if request.Comprehensive {
		md, ok := metadata.FromIncomingContext(ctx)
		var traceID string
		if ok {
			headerValue := md.Get(service.HeaderKey)
			if len(headerValue) > 0 {
				traceID = headerValue[0]
			}
		}
		m.gatherComprehensiveData(&answer, traceID)
	}

	for _, content := range option.Content {
		if content.Greek == request.QuizWord {
			if content.Translation == request.Answer {
				answer.Correct = true
			}
		}
	}

	return &answer, nil
}

// findQuizWord takes a slice and looks for an element in it
func findQuizWord(slice []*pb.Options, val string) bool {
	for _, item := range slice {
		if item.Option == val {
			return true
		}
	}
	return false
}

func (m *MediaServiceImpl) gatherComprehensiveData(answer *pb.ComprehensiveResponse, requestID string) {
	splitID := strings.Split(requestID, "+")

	traceCall := false
	var traceID, parentSpanID string

	if len(splitID) >= 3 {
		traceCall = splitID[2] == "1"
	}

	if len(splitID) >= 1 {
		traceID = splitID[0]
	}
	if len(splitID) >= 2 {
		parentSpanID = splitID[1]
	}

	wordToBeSend := extractBaseWord(answer.QuizWord)

	// Use a WaitGroup to wait for both goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Buffered channels to capture 1 response
	foundInTextChan := make(chan *http.Response, 1)
	similarWordsChan := make(chan *http.Response, 1)
	errChan := make(chan error, 2) // Buffered to hold potential errors from both calls

	go func() {
		defer wg.Done()
		if traceCall {
			herodotosSpan := &pbar.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: parentSpanID,
				SpanId:       comedy.GenerateSpanID(),
				RequestType: &pbar.ParabasisRequest_Span{Span: &pbar.SpanRequest{
					Action: "analyseText",
					Status: fmt.Sprintf("querying Herodotos for word: %s", wordToBeSend),
				}},
			}

			err := m.Streamer.Send(herodotosSpan)
			if err != nil {
				logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
			}
		}
		r := models.AnalyzeTextRequest{Rootword: wordToBeSend}
		jsonBody, err := json.Marshal(r)
		foundInText, err := m.Client.Herodotos().Analyze(jsonBody, requestID)
		if err != nil {
			logging.Error(fmt.Sprintf("could not query any texts for word: %s error: %s", answer.QuizWord, err.Error()))
			errChan <- err
			return
		}
		foundInTextChan <- foundInText
	}()

	go func() {
		defer wg.Done()
		if traceCall {
			alexandrosSpan := &pbar.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: parentSpanID,
				SpanId:       comedy.GenerateSpanID(),
				RequestType: &pbar.ParabasisRequest_Span{Span: &pbar.SpanRequest{
					Action: "analyseText",
					Status: fmt.Sprintf("querying Alexandros for word: %s", wordToBeSend),
				}},
			}

			err := m.Streamer.Send(alexandrosSpan)
			if err != nil {
				logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
			}
		}
		similarWords, err := m.Client.Alexandros().Search(wordToBeSend, "greek", "fuzzy", "false", requestID)
		if err != nil {
			logging.Error(fmt.Sprintf("could not query any similar words for word: %s error: %s", answer.QuizWord, err.Error()))
			errChan <- err
			return
		}
		similarWordsChan <- similarWords
	}()

	// Wait for both goroutines to complete
	wg.Wait()

	// Process responses
	close(errChan)
	close(foundInTextChan)
	close(similarWordsChan)

	for err := range errChan {
		logging.Error(err.Error())
	}

	for foundInText := range foundInTextChan {
		defer foundInText.Body.Close()
		err := json.NewDecoder(foundInText.Body).Decode(&answer.FoundInText)
		if err != nil {
			logging.Error(fmt.Sprintf("error while decoding: %s", err.Error()))
		}
	}

	for similarWords := range similarWordsChan {
		defer similarWords.Body.Close()
		var extended models.ExtendedResponse
		err := json.NewDecoder(similarWords.Body).Decode(&extended)
		if err != nil {
			logging.Error(fmt.Sprintf("error while decoding: %s", err.Error()))
		}

		for _, meros := range extended.Hits {
			answer.SimilarWords = append(answer.SimilarWords, &pb.Meros{
				Greek:      meros.Hit.Greek,
				English:    meros.Hit.English,
				Dutch:      meros.Hit.Dutch,
				LinkedWord: meros.Hit.LinkedWord,
				Original:   meros.Hit.Original,
			})
		}
	}
}

func extractBaseWord(queryWord string) string {
	// Normalize and split the input
	strippedWord := transform.RemoveAccents(strings.ToLower(queryWord))
	splitWord := strings.Split(strippedWord, " ")

	greekPronouns := map[string]bool{"η": true, "ο": true, "το": true}
	cleanWord := func(word string) string {
		return strings.Trim(word, ",.!?-") // Add any other punctuation as needed
	}

	for _, word := range splitWord {
		cleanedWord := cleanWord(word)

		if strings.HasPrefix(cleanedWord, "-") {
			continue
		}

		if _, isPronoun := greekPronouns[cleanedWord]; !isPronoun {
			// If the word is not a pronoun, it's likely the correct word
			return cleanedWord
		}
	}

	return queryWord
}
