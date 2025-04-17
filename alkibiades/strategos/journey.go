package strategos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	pb "github.com/odysseia-greek/apologia/alkibiades/proto"
	"google.golang.org/grpc/metadata"
	"os"
	"time"
)

const (
	OPTIONSEGMENTKEY string = "archytassavedoptions"
	THEME            string = "theme"
	SEGMENT          string = "segment"
)

func (j *JourneyServiceImpl) Health(context.Context, *pb.HealthRequest) (*pb.HealthResponse, error) {
	elasticHealth := j.Elastic.Health().Info()
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
		Version:        os.Getenv("VERSION"),
	}, nil
}

func (j *JourneyServiceImpl) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
	var unparsedResponse []byte
	cacheItem, _ := j.Archytas.Read(OPTIONSEGMENTKEY)
	if cacheItem != nil {
		unparsedResponse = cacheItem
	} else {
		query := quizAggregationQuery()
		logging.Warn(fmt.Sprintf("%v", query))

		elasticResponse, err := j.Elastic.Query().MatchRaw(j.Index, query)
		if err != nil {
			return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
		}

		unparsedResponse = elasticResponse
		err = j.Archytas.Set(OPTIONSEGMENTKEY, string(elasticResponse))
		if err != nil {
			logging.Error(err.Error())
		}
	}

	return parseAggregationResult(unparsedResponse)
}

func (j *JourneyServiceImpl) Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error) {
	var sessionId string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		headerValue := md.Get(config.SessionIdKey)
		if len(headerValue) > 0 {
			sessionId = headerValue[0]
		}
	}

	logging.Debug(fmt.Sprintf("%v", sessionId))

	segmentKey := fmt.Sprintf("%s+%s", request.Theme, request.Segment)
	cacheItem, _ := j.Archytas.Read(segmentKey)

	var option JourneyBasedQuiz

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
		}

		query := j.Elastic.Builder().MultipleMatch(mustQuery)
		elasticResponse, err := j.Elastic.Query().Match(j.Index, query)
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

		err = j.Archytas.Set(segmentKey, string(source))
		if err != nil {
			if err.Error() != "Key not found" {
				logging.Error(fmt.Sprintf("error when writing cache: %s", err.Error()))
			} else {
				logging.Debug(fmt.Sprintf("cache hit: %s", segmentKey))
			}
		}
	}

	quiz := &pb.QuizResponse{
		Theme:       option.Theme,
		Segment:     option.Segment,
		Number:      int32(option.Number),
		Sentence:    option.FullSentence,
		Translation: option.Translation,
		ContextNote: option.ContextNote.Text,
		Quiz:        []*pb.QuizStep{},
	}

	// Add the intro from fixed steps if it exists
	for _, fixedStep := range option.FixedSteps {
		if fixedStep.Type == "intro" {
			quiz.Intro = &pb.Intro{
				Author: fixedStep.Content.Author,

				Work:       fixedStep.Content.Work,
				Background: fixedStep.Content.Background,
			}
			break // We found the intro, no need to continue
		}
	}

	shuffledSteps := j.shuffleAndDistributeSteps(option.RandomSteps)

	// Create quiz steps from random steps
	for _, randomStep := range shuffledSteps {
		var quizStep *pb.QuizStep

		switch randomStep.Type {
		case "match":
			pairs := []*pb.MatchPair{}
			for _, pair := range randomStep.Pairs {
				pairs = append(pairs, &pb.MatchPair{
					Greek:  pair.Greek,
					Answer: pair.Answer,
				})
			}

			for _, pair := range randomStep.Verbs {
				pairs = append(pairs, &pb.MatchPair{
					Greek:  pair.Word,
					Answer: pair.Answer,
				})
			}

			quizStep = &pb.QuizStep{
				Type: &pb.QuizStep_Match{
					Match: &pb.MatchQuiz{
						Instruction: randomStep.Instruction,
						Pairs:       pairs,
					},
				},
			}

		case "trivia":
			options := randomStep.Options
			if len(options) > 4 {
				options = j.limitAndRandomizeOptions(options, randomStep.Answer, 4)
			} else {
				options = j.randomizeOptions(options, randomStep.Answer)
			}

			quizStep = &pb.QuizStep{
				Type: &pb.QuizStep_Trivia{
					Trivia: &pb.TriviaQuiz{
						Question: randomStep.Question,
						Options:  options,
						Answer:   randomStep.Answer,
						Note:     randomStep.Note,
					},
				},
			}

		case "media":
			mediaFiles := []*pb.MediaEntry{}
			for _, media := range randomStep.MediaFiles {
				mediaFiles = append(mediaFiles, &pb.MediaEntry{
					Word:   media.Word,
					Answer: media.Answer,
				})
			}

			quizStep = &pb.QuizStep{
				Type: &pb.QuizStep_Media{
					Media: &pb.MediaDropQuiz{
						Instruction: randomStep.Instruction,
						MediaFiles:  mediaFiles,
					},
				},
			}

		case "structure":
			options := randomStep.Options
			if len(options) > 4 {
				options = j.limitAndRandomizeOptions(options, randomStep.Answer, 4)
			} else {
				options = j.randomizeOptions(options, randomStep.Answer)
			}

			quizStep = &pb.QuizStep{
				Type: &pb.QuizStep_Structure{
					Structure: &pb.StructureQuiz{
						Title:    randomStep.Title,
						Text:     randomStep.Text,
						Question: randomStep.Question,
						Options:  options,
						Answer:   randomStep.Answer,
						Note:     randomStep.NoteOnCorrect,
					},
				},
			}
		}
		if quizStep != nil {
			quiz.Quiz = append(quiz.Quiz, quizStep)
		}
	}

	// Add final translation step
	if option.FinalStep.Type == "translation" {
		options := option.FinalStep.Options
		if len(options) > 4 {
			options = j.limitAndRandomizeOptions(options, option.FinalStep.Answer, 4)
		} else {
			options = j.randomizeOptions(options, option.FinalStep.Answer)
		}

		finalStep := &pb.QuizStep{
			Type: &pb.QuizStep_FinalTranslation{
				FinalTranslation: &pb.TranslationStep{
					Instruction: option.FinalStep.Instruction,
					Options:     options,
					Answer:      option.FinalStep.Answer,
				},
			},
		}
		quiz.Quiz = append(quiz.Quiz, finalStep)
	}

	return quiz, nil
}
