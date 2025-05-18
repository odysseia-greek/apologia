package meletos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/apologia/meletos/model"
	"github.com/odysseia-greek/apologia/meletos/model/queries/multiplechoice"
	"strconv"
)

const (
	MultipleChoiceOptions = "multipleChoiceOptions"
)

func (m *MeletosFixture) iQueryForMultipleChoiceQuizOptions() error {
	query := multiplechoice.Options()
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
			Response model.ThemedOptions `json:"multipleChoiceOptions"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&optionsResponse)

	m.ctx = context.WithValue(m.ctx, MultipleChoiceOptions, optionsResponse.Data.Response)
	return err
}

func (m *MeletosFixture) iUseTheMultipleChoiceOptionsToCreateAQuestion() error {
	options := m.ctx.Value(MultipleChoiceOptions).(model.ThemedOptions)
	randomThemeNumber := m.Randomizer.RandomNumberBaseZero(len(options.Themes))
	randomTheme := options.Themes[randomThemeNumber]

	var randomSet int
	if *randomTheme.MaxSet <= 1 {
		randomSet = 1
	} else {
		randomSet = 1
	}
	randomSetString := strconv.Itoa(randomSet)
	inputBool := true
	doneAfter := int32(2)

	input := model.MultipleQuizInput{
		DoneAfter:       &doneAfter,
		Theme:           randomTheme.Name,
		Set:             &randomSetString,
		ResetProgress:   &inputBool,
		ArchiveProgress: &inputBool,
	}
	query, variables := multiplechoice.Question(input)
	resp, err := m.ForwardGraphql(query, variables)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
	}

	var questionResponse struct {
		Data struct {
			Response model.MultipleChoiceResponse `json:"multipleChoiceQuiz"`
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

func (m *MeletosFixture) iSubmitEachMultipleChoiceOptionOnce() error {
	question := m.ctx.Value(Question).(model.MultipleChoiceResponse)
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
		answer := model.MultipleChoiceAnswerInput{
			Theme:         mediaInput.Theme,
			Set:           mediaInput.Set,
			QuizWord:      question.QuizItem,
			Answer:        option.Option,
			Comprehensive: &comprehensive,
		}

		query, vars := multiplechoice.Answer(answer)
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
				Response model.ComprehensiveResponse `json:"multipleChoiceAnswer"`
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
