package meletos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/apologia/meletos/model"
	"github.com/odysseia-greek/apologia/meletos/model/queries/media"
	"strconv"
)

const (
	MediaOptions = "mediaOptions"
)

func (m *MeletosFixture) iQueryForMediaQuizOptions() error {
	query := media.Options()
	resp, err := m.ForwardGraphql(query, map[string]interface{}{})
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var optionsResponse struct {
		Data struct {
			Response model.AggregatedOptions `json:"mediaOptions"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&optionsResponse)

	m.ctx = context.WithValue(m.ctx, MediaOptions, optionsResponse.Data.Response)

	return err
}

func (m *MeletosFixture) iUseTheMediaOptionsToCreateAQuestion() error {
	options := m.ctx.Value(MediaOptions).(model.AggregatedOptions)
	randomThemeNumber := m.Randomizer.RandomNumberBaseZero(len(options.Themes))
	randomTheme := options.Themes[randomThemeNumber]
	randomSegmentNumber := m.Randomizer.RandomNumberBaseZero(len(randomTheme.Segments))
	randomSegment := randomTheme.Segments[randomSegmentNumber]

	var randomSet int
	if *randomSegment.MaxSet <= 1 {
		randomSet = 1
	} else {
		randomSet = 1
	}
	randomSetString := strconv.Itoa(randomSet)
	inputBool := true
	doneAfter := int32(2)

	input := model.MediaQuizInput{
		DoneAfter:       &doneAfter,
		Theme:           randomTheme.Name,
		Set:             &randomSetString,
		Segment:         randomSegment.Name,
		ResetProgress:   &inputBool,
		ArchiveProgress: &inputBool,
	}
	query, variables := media.Question(input)
	resp, err := m.ForwardGraphql(query, variables)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
	}

	var questionResponse struct {
		Data struct {
			Response model.MediaQuizResponse `json:"mediaQuiz"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&questionResponse)
	if err != nil {
		return err
	}

	m.ctx = context.WithValue(m.ctx, Variables, variables)
	m.ctx = context.WithValue(m.ctx, Question, questionResponse.Data.Response)

	return nil
}

func (m *MeletosFixture) iSubmitEachOptionOnce() error {
	question := m.ctx.Value(Question).(model.MediaQuizResponse)
	variables := m.ctx.Value(Variables).(map[string]interface{})
	var mediaInput model.MediaQuizInput
	err := model.MapToStruct(variables, &mediaInput)
	if err != nil {
		return err
	}

	var counter model.CorrectInCorrect
	var progress *model.ProgressEntry

	comprehensive := false

	for _, option := range question.Options {
		answer := model.MediaAnswerInput{
			Theme:         mediaInput.Theme,
			Set:           mediaInput.Set,
			Segment:       mediaInput.Segment,
			QuizWord:      question.QuizItem,
			Answer:        option.Option,
			Comprehensive: &comprehensive,
			DoneAfter:     mediaInput.DoneAfter,
		}

		query, vars := media.Answer(answer)
		resp, err := m.ForwardGraphql(query, vars)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
		}

		defer resp.Body.Close()

		var answerResponse struct {
			Data struct {
				Response model.ComprehensiveResponse `json:"mediaAnswer"`
			} `json:"data"`
		}

		err = json.NewDecoder(resp.Body).Decode(&answerResponse)
		if err != nil {
			return err
		}

		if *answerResponse.Data.Response.Correct {
			counter.Correct++
		} else {
			counter.Incorrect++
		}

		for _, entry := range answerResponse.Data.Response.Progress {
			if *entry.Greek == *answer.QuizWord {
				progress = entry
				break
			}
		}
	}

	m.ctx = context.WithValue(m.ctx, Responses, counter)
	m.ctx = context.WithValue(m.ctx, Progress, progress)

	return nil
}
