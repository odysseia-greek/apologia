package authorbased

import (
	"github.com/odysseia-greek/apologia/meletos/model"
)

func Options() string {
	return `
		query {
		  authorBasedOptions {
			themes {
					name
					segments {
						name
						maxSet
					}
				}
		  }
		}
	`
}

func Question(question model.AuthorBasedInput) (string, map[string]interface{}) {
	variables, _ := model.StructToMap(question)

	query := `
	query authorBasedQuiz($set: String!, $segment: String!, $theme: String!, $doneAfter: Int!, $resetProgress: Boolean!, $archiveProgress: Boolean!) {
		authorBasedQuiz(
			input: {
				set: $set
				segment: $segment
				theme: $theme
				doneAfter: $doneAfter
				resetProgress: $resetProgress
				archiveProgress: $archiveProgress
			}
		) {
		reference
		translation
		fullSentence
		quiz {
			numberOfItems
			quizItem
			options {
				quizWord
			}
		}
		grammarQuiz {
			correctAnswer
			extraInformation
			wordInText
			options {
				quizWord
			}
		}
		progress {
			greek
			translation
			playCount
			correctCount
			incorrectCount
			lastPlayed
		}
	}
}
`

	return query, variables
}

func Answer(answer model.AuthorBasedAnswerInput) (string, map[string]interface{}) {
	variables, _ := model.StructToMap(answer)

	query := `
	query AuthorBasedAnswer($set: String!, $segment: String!, $theme: String!, $doneAfter: Int!, $answer: String!, $quizWord: String!) {
		authorBasedAnswer(
			input: {
				set: $set
				segment: $segment
				theme: $theme
				doneAfter: $doneAfter
				answer: $answer
				quizWord: $quizWord
			}
		) {
		quizWord
		correct
		finished
		wordsInText
		progress {
			greek
			translation
			playCount
			correctCount
			incorrectCount
			lastPlayed
		}
	}
	}
`
	return query, variables
}

func WordForms(input model.AuthorBasedWordFormsInput) (string, map[string]interface{}) {
	variables, _ := model.StructToMap(input)

	query := `
	query authorBasedWordForms($theme: String!, $segment: String!, $set: String!) {
		authorBasedWordForms(
			input: {
				theme: $theme
				set: $set
				segment: $segment
			}
		) {
			forms {
				dictionaryForm
				wordsInText
			}
		}
	}
	`

	return query, variables
}
