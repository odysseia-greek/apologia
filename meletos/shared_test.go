package meletos

import (
	"fmt"
	"github.com/odysseia-greek/apologia/meletos/model"
	"github.com/stretchr/testify/assert"
)

func (m *MeletosFixture) iShouldHaveIncorrectAndCorrectAnswer(incorrect, correct int) error {
	counter := m.ctx.Value(Responses).(model.CorrectInCorrect)

	err := assertExpectedAndActual(
		assert.Equal, counter.Correct, correct,
		"incorrect number of correct answers",
	)

	if err != nil {
		return err
	}

	err = assertExpectedAndActual(
		assert.Equal, counter.Incorrect, incorrect,
		"incorrect number of incorrect answers",
	)

	return err
}

func (m *MeletosFixture) theProgressShouldBeIncorrectAndCorrectAnswer(incorrect, correct int) error {
	progressEntry := m.ctx.Value(Progress).(*model.ProgressEntry)

	if int32(incorrect) == *progressEntry.IncorrectCount && int32(correct) == *progressEntry.CorrectCount {
		return nil
	}

	return fmt.Errorf("no progress found for incorrect: %d and correct: %d", incorrect, correct)
}
