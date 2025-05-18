package meletos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/apologia/meletos/model"
	"github.com/odysseia-greek/apologia/meletos/model/queries/journey"
	"io"
	"log"
)

const (
	JourneyOptions = "journeyOptions"
)

type RawJourneySegmentQuiz struct {
	Theme       string            `json:"theme"`
	Segment     string            `json:"segment"`
	Number      int32             `json:"number"`
	Sentence    string            `json:"sentence"`
	Translation string            `json:"translation"`
	ContextNote *string           `json:"contextNote"`
	Intro       *model.QuizIntro  `json:"intro"`
	Quiz        []json.RawMessage `json:"quiz"`
}

func (m *MeletosFixture) aNewJourneyIsReturnedWithATranslationAndSentence() error {
	question := m.ctx.Value(Question).(*model.JourneySegmentQuiz)
	if question.Sentence == "" || question.Translation == "" {
		return fmt.Errorf("journey has no greek or english sentence")
	}

	return nil
}

func (m *MeletosFixture) aShortBackgroundOnTheTextShouldExist() error {
	question := m.ctx.Value(Question).(*model.JourneySegmentQuiz)
	if question.Intro.Work == "" || question.Intro.Author == "" || question.Intro.Background == "" {
		return fmt.Errorf("journey has no intro")
	}

	return nil
}

func (m *MeletosFixture) iQueryForJourneyQuizOptions() error {
	query := journey.Options()
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
			Response model.JourneyOptions `json:"journeyOptions"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&optionsResponse)

	m.ctx = context.WithValue(m.ctx, JourneyOptions, optionsResponse.Data.Response)
	return err
}

func (m *MeletosFixture) iUseTheJourneyOptionsToCreateAQuestion() error {
	options := m.ctx.Value(JourneyOptions).(model.JourneyOptions)
	randomThemeNumber := m.Randomizer.RandomNumberBaseZero(len(options.Themes))
	randomTheme := options.Themes[randomThemeNumber]
	randomSegmentNumber := m.Randomizer.RandomNumberBaseZero(len(randomTheme.Segments))
	randomSegment := randomTheme.Segments[randomSegmentNumber]

	input := model.JourneyQuizInput{
		Theme:   randomTheme.Name,
		Segment: randomSegment.Name,
	}

	query, variables := journey.Question(input)
	resp, err := m.ForwardGraphql(query, variables)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read Sokrates response: %w", err)
	}

	var raw struct {
		Data struct {
			Section RawJourneySegmentQuiz `json:"journeyQuiz"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}

	var quizSections []model.QuizSection

	for _, rawItem := range raw.Data.Section.Quiz {
		type quizProbe struct {
			Typename string `json:"__typename"`
		}

		var probe quizProbe
		if err := json.Unmarshal(rawItem, &probe); err != nil {
			log.Println("Unmarshal probe failed:", err)
			continue
		}

		switch probe.Typename {
		case "StructureQuiz":
			var s model.StructureQuiz
			if err := json.Unmarshal(rawItem, &s); err == nil {
				quizSections = append(quizSections, &s)
			}
		case "MatchQuiz":
			var m model.MatchQuiz
			if err := json.Unmarshal(rawItem, &m); err == nil {
				quizSections = append(quizSections, &m)
			}
		case "TriviaQuiz":
			var t model.TriviaQuiz
			if err := json.Unmarshal(rawItem, &t); err == nil {
				quizSections = append(quizSections, &t)
			}
		case "MediaQuiz":
			var t model.MediaQuiz
			if err := json.Unmarshal(rawItem, &t); err == nil {
				quizSections = append(quizSections, &t)
			}
		case "FinalTranslationQuiz":
			var t model.FinalTranslationQuiz
			if err := json.Unmarshal(rawItem, &t); err == nil {
				quizSections = append(quizSections, &t)
			}
		default:
			log.Println("Unknown quiz section type:", probe.Typename)
		}
	}

	questionResponse := &model.JourneySegmentQuiz{
		Theme:       raw.Data.Section.Theme,
		Segment:     raw.Data.Section.Segment,
		Number:      raw.Data.Section.Number,
		Sentence:    raw.Data.Section.Sentence,
		Translation: raw.Data.Section.Translation,
		ContextNote: raw.Data.Section.ContextNote,
		Intro:       raw.Data.Section.Intro,
		Quiz:        quizSections,
	}

	m.ctx = context.WithValue(m.ctx, Question, questionResponse)

	return nil
}

func (m *MeletosFixture) theQuizHasDifferentTypesOfQuestionsEmbedded() error {
	question := m.ctx.Value(Question).(*model.JourneySegmentQuiz)

	if len(question.Quiz) < 2 {
		return fmt.Errorf("journey has no quiz")
	}

	return nil
}
