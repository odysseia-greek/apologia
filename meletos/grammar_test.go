package meletos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/apologia/meletos/model"
	"github.com/odysseia-greek/apologia/meletos/model/queries/grammar"
	"strconv"
)

const (
	GrammarOptions = "grammarOptions"
)

func (m *MeletosFixture) iQueryForGrammarQuizOptions() error {
	query := grammar.Options()
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
			Response model.GrammarOptions `json:"grammarOptions"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&optionsResponse)

	m.ctx = context.WithValue(m.ctx, GrammarOptions, optionsResponse.Data.Response)
	return err
}

func (m *MeletosFixture) iUseTheGrammarOptionsToCreateAQuestion() error {
	options := m.ctx.Value(GrammarOptions).(model.GrammarOptions)
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

	input := model.GrammarQuizInput{
		DoneAfter:       &doneAfter,
		Theme:           randomTheme.Name,
		Set:             &randomSetString,
		Segment:         randomSegment.Name,
		ResetProgress:   &inputBool,
		ArchiveProgress: &inputBool,
	}
	query, variables := grammar.Question(input)
	resp, err := m.ForwardGraphql(query, variables)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
	}

	var questionResponse struct {
		Data struct {
			Response model.GrammarQuizResponse `json:"grammarQuiz"`
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

func (m *MeletosFixture) iSubmitEachGrammarOptionOnce() error {
	question := m.ctx.Value(Question).(model.GrammarQuizResponse)
	variables := m.ctx.Value(Variables).(map[string]interface{})
	var grammarInput model.GrammarQuizInput
	err := model.MapToStruct(variables, &grammarInput)
	if err != nil {
		return err
	}

	var counter model.CorrectInCorrect
	var progress *model.ProgressEntry

	comprehensive := false

	for _, option := range question.Options {
		answer := model.GrammarAnswerInput{
			Theme:          grammarInput.Theme,
			Set:            grammarInput.Set,
			Segment:        grammarInput.Segment,
			QuizWord:       question.QuizItem,
			Answer:         option.Option,
			DoneAfter:      grammarInput.DoneAfter,
			DictionaryForm: question.DictionaryForm,
			Comprehensive:  &comprehensive,
		}

		query, vars := grammar.Answer(answer)
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
				Response model.ComprehensiveResponse `json:"grammarAnswer"`
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
