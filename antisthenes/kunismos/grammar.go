package kunismos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/agora/plato/transform"
	pb "github.com/odysseia-greek/apologia/antisthenes/proto"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"math/rand/v2"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	OPTIONSEGMENTKEY string = "archytassavedoptions"
	THEME            string = "theme"
	SET              string = "set"
	SEGMENT          string = "segment"
)

func (g *GrammarServiceImpl) Health(context.Context, *pb.HealthRequest) (*pb.HealthResponse, error) {
	elasticHealth := g.Elastic.Health().Info()
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
		Version:        g.Version,
	}, nil
}

func (g *GrammarServiceImpl) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
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

func (g *GrammarServiceImpl) Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error) {
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

	quiz := &pb.QuizResponse{
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
	quiz.Options = append(quiz.Options, &pb.GrammarOptions{
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
			quiz.Options = append(quiz.Options, &pb.GrammarOptions{
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
			quiz.Progress = append(quiz.Progress, &pb.ProgressEntry{
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

func (g *GrammarServiceImpl) Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.ComprehensiveResponse, error) {
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

	answer := pb.ComprehensiveResponse{Correct: false, QuizWord: request.DictionaryForm}

	if request.Comprehensive {
		md, ok := metadata.FromIncomingContext(ctx)
		var traceID string
		if ok {
			headerValue := md.Get(service.HeaderKey)
			if len(headerValue) > 0 {
				traceID = headerValue[0]
			}
		}
		g.gatherComprehensiveData(&answer, traceID)
	}

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
			answer.Progress = append(answer.Progress, &pb.ProgressEntry{
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

func (g *GrammarServiceImpl) gatherComprehensiveData(answer *pb.ComprehensiveResponse, requestID string) {
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

			err := g.Streamer.Send(herodotosSpan)
			if err != nil {
				logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
			}
		}
		r := models.AnalyzeTextRequest{Rootword: wordToBeSend}
		jsonBody, err := json.Marshal(r)
		foundInText, err := g.Client.Herodotos().Analyze(jsonBody, requestID)
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

			err := g.Streamer.Send(alexandrosSpan)
			if err != nil {
				logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
			}
		}
		similarWords, err := g.Client.Alexandros().Search(wordToBeSend, "greek", "fuzzy", "false", requestID)
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
		var foundInTextModel models.AnalyzeTextResponse
		err := json.NewDecoder(foundInText.Body).Decode(&foundInTextModel)
		if err != nil {
			logging.Error(fmt.Sprintf("error while decoding: %s", err.Error()))
		}

		grpcModel := &pb.AnalyzeTextResponse{
			Rootword:     foundInTextModel.Rootword,
			PartOfSpeech: foundInTextModel.PartOfSpeech,
		}

		var conj []*pb.Conjugations
		for _, conjugation := range foundInTextModel.Conjugations {
			conj = append(conj, &pb.Conjugations{
				Word: conjugation.Word,
				Rule: conjugation.Rule,
			})
		}

		grpcModel.Conjugations = conj

		var result []*pb.AnalyzeResult
		for _, text := range foundInTextModel.Results {
			result = append(result, &pb.AnalyzeResult{
				ReferenceLink: text.ReferenceLink,
				Author:        text.Author,
				Book:          text.Book,
				Reference:     text.Reference,
				Text: &pb.Rhema{
					Greek:        text.Text.Greek,
					Translations: text.Text.Translations,
					Section:      text.Text.Section,
				},
			})
		}

		grpcModel.Texts = result

		answer.FoundInText = grpcModel
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

func sliceToSet(words []string) map[string]struct{} {
	set := make(map[string]struct{}, len(words))
	for _, w := range words {
		set[w] = struct{}{}
	}
	return set
}

// findQuizWord takes a slice and looks for an element in it
func findQuizWord(slice []*pb.GrammarOptions, val string) bool {
	for _, item := range slice {
		if item.Option == val {
			return true
		}
	}
	return false
}
