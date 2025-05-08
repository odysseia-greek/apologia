package multiplechoice

import (
	"github.com/odysseia-greek/apologia/meletos/model"
)

func Options() string {
	return `
		query {
		  multipleChoiceOptions {
			themes {
			  name
              maxSet
			}
		  }
		}
	`
}

func Answer(answer model.MultipleChoiceAnswerInput) (string, map[string]interface{}) {
	variables, _ := model.StructToMap(answer)

	query := `query multipleChoiceAnswer($set: String!, $theme: String!, $quizWord: String!, $answer: String!, $comprehensive: Boolean!) {
		multipleChoiceAnswer(
			input: {
				set: $set
				theme: $theme
				quizWord: $quizWord	
				answer: $answer
				comprehensive: $comprehensive
			}
		) {
			correct
			finished
			quizWord
			similarWords {
			  greek
			  english
			}
			foundInText {
			  rootword
			  conjugations {
				word
				rule
			  }
			  texts {
				author
				book
				text {
				  translations
				  greek
				}
			  }
			}
				progress {
					greek
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

func Question(question model.MultipleQuizInput) (string, map[string]interface{}) {
	variables, _ := model.StructToMap(question)

	query := `
	query multipleChoiceQuiz($set: String!, $theme: String!, $doneAfter: Int!, $resetProgress: Boolean!, $archiveProgress: Boolean!) {
		multipleChoiceQuiz(
			input: {
				set: $set
				theme: $theme
				doneAfter: $doneAfter
				resetProgress: $resetProgress
				archiveProgress: $archiveProgress
			}
		) {
			numberOfItems
			quizItem
			options {
				option
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
