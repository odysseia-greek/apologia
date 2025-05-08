package meletos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/apologia/meletos/model"
	"github.com/odysseia-greek/apologia/meletos/model/queries/authorbased"
	"strconv"
)

const (
	AuthorbasedOptions = "authorOptions"
	WordForms          = "wordForms"
)

func (m *MeletosFixture) iQueryForAuthorbasedQuizOptions() error {
	query := authorbased.Options()
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
			Response model.AggregatedOptions `json:"authorBasedOptions"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&optionsResponse)

	m.ctx = context.WithValue(m.ctx, AuthorbasedOptions, optionsResponse.Data.Response)
	return err
}

func (m *MeletosFixture) iSubmitEachAuthorbasedOptionOnce() error {
	question := m.ctx.Value(Question).(model.AuthorBasedResponse)
	variables := m.ctx.Value(Variables).(map[string]interface{})
	var authorbasedInput model.AuthorBasedInput
	err := model.MapToStruct(variables, &authorbasedInput)
	if err != nil {
		return err
	}

	var counter model.CorrectInCorrect
	var progress *model.ProgressEntry

	for _, option := range question.Quiz.Options {
		answer := model.AuthorBasedAnswerInput{
			Theme:     authorbasedInput.Theme,
			Set:       authorbasedInput.Set,
			Segment:   authorbasedInput.Segment,
			QuizWord:  question.Quiz.QuizItem,
			Answer:    option.QuizWord,
			DoneAfter: authorbasedInput.DoneAfter,
		}

		query, vars := authorbased.Answer(answer)
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
				Response model.AuthorBasedAnswerResponse `json:"authorBasedAnswer"`
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

func (m *MeletosFixture) iUseTheAuthorbasedOptionsToCreateAQuestion() error {
	options := m.ctx.Value(AuthorbasedOptions).(model.AggregatedOptions)
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

	input := model.AuthorBasedInput{
		DoneAfter:       &doneAfter,
		Theme:           randomTheme.Name,
		Set:             &randomSetString,
		Segment:         randomSegment.Name,
		ResetProgress:   &inputBool,
		ArchiveProgress: &inputBool,
	}
	query, variables := authorbased.Question(input)
	resp, err := m.ForwardGraphql(query, variables)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
	}

	var questionResponse struct {
		Data struct {
			Response model.AuthorBasedResponse `json:"authorBasedQuiz"`
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

func (m *MeletosFixture) grammarOptionsShouldBeEmbeddedIntoTheQuizForSomeWords() error {
	grammarQuiz := m.ctx.Value(GrammarQuiz).([]model.GrammarQuizAdded)
	if len(grammarQuiz) == 0 {
		return fmt.Errorf("no grammar quiz was created")
	}

	return nil
}

func (m *MeletosFixture) iCreateAQuizThatHasTheName(name string) error {
	options := m.ctx.Value(AuthorbasedOptions).(model.AggregatedOptions)
	var theme model.Theme
	for _, t := range options.Themes {
		if *t.Name == name {
			theme = *t
			break
		}
	}
	randomSegmentNumber := m.Randomizer.RandomNumberBaseZero(len(theme.Segments))
	randomSegment := theme.Segments[randomSegmentNumber]

	var randomSet int
	if *randomSegment.MaxSet <= 1 {
		randomSet = 1
	} else {
		randomSet = 1
	}
	randomSetString := strconv.Itoa(randomSet)
	inputBool := true
	doneAfter := int32(2)

	input := model.AuthorBasedInput{
		DoneAfter:       &doneAfter,
		Theme:           theme.Name,
		Set:             &randomSetString,
		Segment:         randomSegment.Name,
		ResetProgress:   &inputBool,
		ArchiveProgress: &inputBool,
	}

	var grammarQuiz []model.GrammarQuizAdded
	var found bool
	for !found {
		query, variables := authorbased.Question(input)
		resp, err := m.ForwardGraphql(query, variables)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
		}

		defer resp.Body.Close()

		var questionResponse struct {
			Data struct {
				Response model.AuthorBasedResponse `json:"authorBasedQuiz"`
			} `json:"data"`
		}

		err = json.NewDecoder(resp.Body).Decode(&questionResponse)
		if err != nil {
			return err
		}

		if questionResponse.Data.Response.GrammarQuiz != nil {
			for _, q := range questionResponse.Data.Response.GrammarQuiz {
				grammarQuiz = append(grammarQuiz, *q)
			}
			found = true
		}
	}

	m.ctx = context.WithValue(m.ctx, GrammarQuiz, grammarQuiz)

	return nil
}

func (m *MeletosFixture) thatQuestionHasAGreekEnglishSentence() error {
	question := m.ctx.Value(Question).(model.AuthorBasedResponse)
	if *question.FullSentence == "" || *question.Translation == "" {
		return fmt.Errorf("question has no greek or english sentence")
	}

	return nil
}

func (m *MeletosFixture) thatQuestionHasAReferenceToTheTextModule() error {
	question := m.ctx.Value(Question).(model.AuthorBasedResponse)
	if *question.Reference == "" {
		return fmt.Errorf("question has no reference to the text module")
	}

	return nil
}

func (m *MeletosFixture) iQueryTheWordFormsForASegment() error {
	options := m.ctx.Value(AuthorbasedOptions).(model.AggregatedOptions)
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

	input := model.AuthorBasedWordFormsInput{
		Theme:   randomTheme.Name,
		Set:     &randomSetString,
		Segment: randomSegment.Name,
	}
	query, variables := authorbased.WordForms(input)
	resp, err := m.ForwardGraphql(query, variables)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
	}

	var formsResponse struct {
		Data struct {
			Response model.AuthorBasedWordFormsResponse `json:"authorBasedWordForms"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&formsResponse)
	if err != nil {
		return err
	}

	m.ctx = context.WithValue(m.ctx, WordForms, formsResponse.Data.Response)

	return nil
}

func (m *MeletosFixture) theWordsShouldBeReturnedAsTheyAppearInTheText() error {
	wordForms := m.ctx.Value(WordForms).(model.AuthorBasedWordFormsResponse)
	if len(wordForms.Forms) < 1 {
		return fmt.Errorf("no word forms were returned")
	}

	return nil
}
