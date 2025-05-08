package meletos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/apologia/meletos/model"
	"github.com/odysseia-greek/apologia/meletos/model/queries/dialogue"
	"strconv"
)

const (
	DialogueOptions = "dialogueOptions"
)

func (m *MeletosFixture) iQueryForDialogueQuizOptions() error {
	query := dialogue.Options()
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
			Response model.ThemedOptions `json:"dialogueOptions"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&optionsResponse)

	m.ctx = context.WithValue(m.ctx, DialogueOptions, optionsResponse.Data.Response)
	return err
}

func (m *MeletosFixture) iUseTheDialogueOptionsToCreateAQuestion() error {
	options := m.ctx.Value(DialogueOptions).(model.ThemedOptions)
	randomThemeNumber := m.Randomizer.RandomNumberBaseZero(len(options.Themes))
	randomTheme := options.Themes[randomThemeNumber]

	var randomSet int
	if *randomTheme.MaxSet <= 1 {
		randomSet = 1
	} else {
		randomSet = 1
	}
	randomSetString := strconv.Itoa(randomSet)

	input := model.DialogueQuizInput{
		Theme: randomTheme.Name,
		Set:   &randomSetString,
	}
	query, variables := dialogue.Question(input)
	resp, err := m.ForwardGraphql(query, variables)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
	}

	var questionResponse struct {
		Data struct {
			Response model.DialogueQuizResponse `json:"dialogueQuiz"`
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

func (m *MeletosFixture) thePercentageShouldBe(percentage int) error {
	answer := m.ctx.Value(Responses).(model.DialogueAnswer)

	if answer.Percentage == nil {
		return fmt.Errorf("percentage value is nil")
	}

	parsedPercentage := int(*answer.Percentage)
	if parsedPercentage != percentage {
		return fmt.Errorf("percentage should be the same as %d, but was %d", percentage, parsedPercentage)
	}
	return nil
}

func (m *MeletosFixture) thePercentageShouldBeLowerThan(percentage int) error {
	answer := m.ctx.Value(Responses).(model.DialogueAnswer)

	if answer.Percentage == nil {
		return fmt.Errorf("percentage value is nil")
	}

	parsedPercentage := int(*answer.Percentage)
	if parsedPercentage > percentage {
		return fmt.Errorf("percentage should be lower than %d, but was %d", percentage, parsedPercentage)
	}
	return nil
}

func (m *MeletosFixture) wronglyPlacedShouldBeEmpty() error {
	answer := m.ctx.Value(Responses).(model.DialogueAnswer)
	if answer.WronglyPlaced != nil {
		return fmt.Errorf("wronglyPlaced should be nil")
	}

	return nil
}

func (m *MeletosFixture) wronglyPlacedShouldHoldAReferenceToTheCorrectPlace() error {
	answer := m.ctx.Value(Responses).(model.DialogueAnswer)
	if answer.WronglyPlaced == nil {
		return fmt.Errorf("wronglyPlaced value is nil")
	}

	if len(answer.WronglyPlaced) == 0 {
		return fmt.Errorf("wronglyPlaced should not be empty")
	}

	for _, item := range answer.WronglyPlaced {
		if item.Place == item.CorrectPlace {
			return fmt.Errorf("place value should be different from correctPlace value")
		}
	}

	return nil
}

func (m *MeletosFixture) iSubmitWithAPerfectInput() error {
	question := m.ctx.Value(Question).(model.DialogueQuizResponse)
	variables := m.ctx.Value(Variables).(map[string]interface{})
	var dialogueInput model.DialogueAnswerInput
	err := model.MapToStruct(variables, &dialogueInput)
	if err != nil {
		return err
	}

	// Create a deep copy of the content
	var content []*model.DialogueInputContent
	for _, cont := range question.Content {
		contentCopy := &model.DialogueInputContent{
			Translation: cont.Translation,
			Greek:       cont.Greek,
			Place:       cont.Place,
			Speaker:     cont.Speaker,
		}

		content = append(content, contentCopy)
	}

	// Set the answer in dialogueInput
	dialogueInput.Content = content

	// Put dialogueInput back into variables
	newVariables, err := model.StructToMap(dialogueInput)
	if err != nil {
		return err
	}

	// Submit the answer
	query, _ := dialogue.Answer(dialogueInput)
	resp, err := m.ForwardGraphql(query, newVariables)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
	}

	var answerResponse struct {
		Data struct {
			Response model.DialogueAnswer `json:"dialogueAnswer"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&answerResponse)
	if err != nil {
		return err
	}

	m.ctx = context.WithValue(m.ctx, Responses, answerResponse.Data.Response)

	return nil
}

func (m *MeletosFixture) iSubmitWithAtLeastOneSectionWronlyPlaced() error {
	question := m.ctx.Value(Question).(model.DialogueQuizResponse)
	variables := m.ctx.Value(Variables).(map[string]interface{})
	var dialogueInput model.DialogueAnswerInput
	err := model.MapToStruct(variables, &dialogueInput)
	if err != nil {
		return err
	}

	// Create a deep copy of the content
	var content []*model.DialogueInputContent
	for _, cont := range question.Content {
		contentCopy := &model.DialogueInputContent{
			Translation: cont.Translation,
			Greek:       cont.Greek,
			Place:       cont.Place,
			Speaker:     cont.Speaker,
		}

		content = append(content, contentCopy)
	}

	// Choose two random positions to swap
	length := len(content)
	if length < 2 {
		return fmt.Errorf("not enough content items to swap")
	}

	pos1 := m.Randomizer.RandomNumberBaseZero(length)
	pos2 := m.Randomizer.RandomNumberBaseZero(length)

	// Make sure we're not swapping the same position
	for pos1 == pos2 {
		pos2 = m.Randomizer.RandomNumberBaseZero(length)
	}

	// Swap the Place values
	place1 := content[pos1].Place
	content[pos1].Place = content[pos2].Place
	content[pos2].Place = place1

	// Set the answer in dialogueInput
	dialogueInput.Content = content

	// Put dialogueInput back into variables
	newVariables, err := model.StructToMap(dialogueInput)
	if err != nil {
		return err
	}

	// Submit the answer
	query, _ := dialogue.Answer(dialogueInput)
	resp, err := m.ForwardGraphql(query, newVariables)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
	}

	var answerResponse struct {
		Data struct {
			Response model.DialogueAnswer `json:"dialogueAnswer"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&answerResponse)
	if err != nil {
		return err
	}

	m.ctx = context.WithValue(m.ctx, Responses, answerResponse.Data.Response)

	return nil
}
