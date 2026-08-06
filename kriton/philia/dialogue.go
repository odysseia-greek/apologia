package philia

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	koinosv1 "github.com/odysseia-greek/apologia/diotima/gen/go/koinos/v1"
	v1 "github.com/odysseia-greek/apologia/kriton/gen/go/v1"
	"github.com/odysseia-greek/attike/aristophanes/comedy"

	"strconv"
	"time"
)

const (
	THEME            string = "theme"
	SET              string = "set"
	OPTIONSEGMENTKEY string = "archytassavedoptions"
)

func (d *DialogueServiceImpl) Health(context.Context, *koinosv1.HealthRequest) (*koinosv1.HealthResponse, error) {
	elasticHealth := d.Elastic.Health().Info()
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
		Version:        d.Version,
	}, nil
}

func (d *DialogueServiceImpl) Options(ctx context.Context, request *koinosv1.OptionsRequest) (*v1.AggregatedOptions, error) {
	var unparsedResponse []byte
	cacheItem, _ := d.Archytas.Get(OPTIONSEGMENTKEY)
	if cacheItem != nil {
		unparsedResponse = cacheItem
	} else {
		query := quizAggregationQuery()

		elasticResponse, err := d.Elastic.Query().MatchRawWithContext(ctx, d.Index, query)
		if err != nil {
			return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
		}

		unparsedResponse = elasticResponse
		err = d.Archytas.SetBytes(OPTIONSEGMENTKEY, elasticResponse)
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

func (d *DialogueServiceImpl) Question(ctx context.Context, request *v1.CreationRequest) (*v1.QuizResponse, error) {
	segmentKey := fmt.Sprintf("%s+%s", request.Theme, request.Set)
	cacheItem, _ := d.Archytas.Get(segmentKey)

	var quiz models.DialogueQuiz

	if cacheItem != nil {
		err := json.Unmarshal(cacheItem, &quiz)
		if err != nil {
			return nil, err
		}

		go comedy.CacheSpan(string(cacheItem), segmentKey, ctx, d.Streamer)
	} else {
		mustQuery := []map[string]string{
			{
				THEME: request.Theme,
			},
			{
				SET: request.Set,
			},
		}

		query := d.Elastic.Builder().MultipleMatch(mustQuery)
		elasticResponse, err := d.Elastic.Query().MatchWithContext(ctx, d.Index, query)
		if err != nil {
			return nil, err
		}
		if len(elasticResponse.Hits.Hits) == 0 {
			return nil, fmt.Errorf("no hits found in Elastic")
		}

		go comedy.DatabaseSpan(query, elasticResponse.Hits.Total.Value, elasticResponse.Took, ctx, d.Streamer)

		source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
		err = json.Unmarshal(source, &quiz)
		if err != nil {
			return nil, err
		}
	}

	result := &v1.QuizResponse{
		QuizMetadata: &v1.QuizMetadata{
			Language: quiz.QuizMetadata.Language,
		},
		Theme:     quiz.Theme,
		Set:       strconv.Itoa(quiz.Set),
		Segment:   quiz.Dialogue.Section,
		Reference: quiz.Reference,
	}

	dialogue := &v1.Dialogue{
		Introduction:  quiz.Dialogue.Introduction,
		Section:       quiz.Dialogue.Section,
		LinkToPerseus: quiz.Dialogue.LinkToPerseus,
	}

	for _, speaker := range quiz.Dialogue.Speakers {
		dialogue.Speakers = append(dialogue.Speakers, &v1.Speaker{
			Name:        speaker.Name,
			Shorthand:   speaker.Shorthand,
			Translation: speaker.Translation,
		})
	}

	result.Dialogue = dialogue

	for _, content := range quiz.Content {
		dialogueContent := &v1.DialogueContent{
			Translation: content.Translation,
			Greek:       content.Greek,
			Place:       int32(content.Place),
			Speaker:     content.Speaker,
		}

		result.Content = append(result.Content, dialogueContent)
	}

	return result, nil
}

func (d *DialogueServiceImpl) Answer(ctx context.Context, request *v1.AnswerRequest) (*v1.AnswerResponse, error) {
	segmentKey := fmt.Sprintf("%s+%s", request.Theme, request.Set)
	cacheItem, _ := d.Archytas.Get(segmentKey)

	var option models.DialogueQuiz

	if cacheItem != nil {
		err := json.Unmarshal(cacheItem, &option)
		if err != nil {
			return nil, err
		}

		go comedy.CacheSpan(string(cacheItem), segmentKey, ctx, d.Streamer)
	} else {
		mustQuery := []map[string]string{
			{
				THEME: request.Theme,
			},
			{
				SET: request.Set,
			},
		}

		query := d.Elastic.Builder().MultipleMatch(mustQuery)
		elasticResponse, err := d.Elastic.Query().MatchWithContext(ctx, d.Index, query)
		if err != nil {
			return nil, err
		}
		if len(elasticResponse.Hits.Hits) == 0 {
			return nil, fmt.Errorf("no hits found in Elastic")
		}

		go comedy.DatabaseSpan(query, elasticResponse.Hits.Total.Value, elasticResponse.Took, ctx, d.Streamer)

		source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
		err = json.Unmarshal(source, &option)
		if err != nil {
			return nil, err
		}
	}

	answer := &v1.AnswerResponse{
		Percentage:    0,
		Input:         request.Content,
		Answer:        []*v1.DialogueContent{},
		WronglyPlaced: nil,
	}

	for _, content := range option.Content {
		answer.Answer = append(answer.Answer, &v1.DialogueContent{
			Translation: content.Translation,
			Greek:       content.Greek,
			Place:       int32(content.Place),
			Speaker:     content.Speaker,
		})
	}

	var correctPlace int
	var wrongPlace int

	for _, dialogue := range request.Content {
		verifiedContent := option.Content[dialogue.Place-1]
		if verifiedContent.Greek == dialogue.Greek && int32(verifiedContent.Place) == dialogue.Place {
			correctPlace++
		} else {
			correctedPlacing := &v1.DialogueCorrection{
				Translation:  dialogue.Translation,
				Greek:        dialogue.Greek,
				Place:        dialogue.Place,
				Speaker:      dialogue.Speaker,
				CorrectPlace: 0,
			}

			for _, corrected := range option.Content {
				if corrected.Greek == dialogue.Greek && corrected.Speaker == dialogue.Speaker {
					correctedPlacing.CorrectPlace = int32(corrected.Place)
				}
			}

			answer.WronglyPlaced = append(answer.WronglyPlaced, correctedPlacing)
			wrongPlace++
		}
	}

	total := correctPlace + wrongPlace
	totalProgress := 0.0
	if total > 0 {
		totalProgress = float64(correctPlace) / float64(total) * 100
	}

	answer.Percentage = totalProgress

	return answer, nil
}
