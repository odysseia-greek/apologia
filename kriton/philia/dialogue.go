package philia

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/odysseia-greek/agora/plato/models"
	pb "github.com/odysseia-greek/apologia/kriton/proto"
	"strconv"
	"time"
)

const (
	THEME string = "theme"
	SET   string = "set"
)

func (d *DialogueServiceImpl) Health(context.Context, *pb.HealthRequest) (*pb.HealthResponse, error) {
	elasticHealth := d.Elastic.Health().Info()
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

func (d *DialogueServiceImpl) Options(ctx context.Context, request *pb.OptionsRequest) (*pb.AggregatedOptions, error) {
	query := quizAggregationQuery()

	elasticResult, err := d.Elastic.Query().MatchRaw(d.Index, query)
	if err != nil {
		return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
	}

	result, err := parseAggregationResult(elasticResult)
	if err != nil {
		return nil, fmt.Errorf("error in elasticSearch: %s", err.Error())
	}

	return result, nil
}

func (d *DialogueServiceImpl) Question(ctx context.Context, request *pb.CreationRequest) (*pb.QuizResponse, error) {
	mustQuery := []map[string]string{
		{
			THEME: request.Theme,
		},
		{
			SET: request.Set,
		},
	}
	query := d.Elastic.Builder().MultipleMatch(mustQuery)
	elasticResponse, err := d.Elastic.Query().Match(d.Index, query)
	if err != nil {
		return nil, err
	}

	if elasticResponse.Hits.Hits == nil || len(elasticResponse.Hits.Hits) == 0 {
		return nil, errors.New("no hits found in query")
	}

	var quiz models.DialogueQuiz

	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &quiz)
	if err != nil {
		return nil, err
	}

	result := &pb.QuizResponse{
		QuizMetadata: &pb.QuizMetadata{
			Language: quiz.QuizMetadata.Language,
		},
		Theme:     quiz.Theme,
		Set:       strconv.Itoa(quiz.Set),
		Segment:   quiz.Segment,
		Reference: quiz.Reference,
	}

	dialogue := &pb.Dialogue{
		Introduction:  quiz.Dialogue.Introduction,
		Section:       quiz.Dialogue.Section,
		LinkToPerseus: quiz.Dialogue.LinkToPerseus,
	}

	for _, speaker := range quiz.Dialogue.Speakers {
		dialogue.Speakers = append(dialogue.Speakers, &pb.Speaker{
			Name:        speaker.Name,
			Shorthand:   speaker.Shorthand,
			Translation: speaker.Translation,
		})
	}

	result.Dialogue = dialogue

	for _, content := range quiz.Content {
		dialogueContent := &pb.DialogueContent{
			Translation: content.Translation,
			Greek:       content.Greek,
			Place:       int32(content.Place),
			Speaker:     content.Speaker,
		}

		result.Content = append(result.Content, dialogueContent)
	}

	return result, nil
}

func (d *DialogueServiceImpl) Answer(ctx context.Context, request *pb.AnswerRequest) (*pb.AnswerResponse, error) {
	mustQuery := []map[string]string{
		{
			THEME: request.Theme,
		},
		{
			SET: request.Set,
		},
	}

	query := d.Elastic.Builder().MultipleMatch(mustQuery)
	elasticResponse, err := d.Elastic.Query().Match(d.Index, query)
	if err != nil {
		return nil, err
	}

	if len(elasticResponse.Hits.Hits) == 0 {
		return nil, fmt.Errorf("no hits found in Elastic")
	}

	var option models.DialogueQuiz
	source, _ := json.Marshal(elasticResponse.Hits.Hits[0].Source)
	err = json.Unmarshal(source, &option)
	if err != nil {
		return nil, err
	}

	answer := &pb.AnswerResponse{
		Percentage:    0,
		Input:         request.Content,
		Answer:        []*pb.DialogueContent{},
		WronglyPlaced: nil,
	}

	for _, content := range option.Content {
		answer.Answer = append(answer.Answer, &pb.DialogueContent{
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
			correctedPlacing := &pb.DialogueCorrection{
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
