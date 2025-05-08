package media

import (
	"github.com/odysseia-greek/apologia/meletos/model"
)

func Options() string {
	return `
		query {
		  mediaOptions {
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

func Answer(answer model.MediaAnswerInput) (string, map[string]interface{}) {
	variables, _ := model.StructToMap(answer)

	query := `
	query MediaAnswer($set: String!, $segment: String!, $theme: String!, $quizWord: String!, $answer: String!, $doneAfter: Int!, $comprehensive: Boolean!) {
		mediaAnswer(
			input: {
				set: $set
				segment: $segment
				theme: $theme
				quizWord: $quizWord	
				answer: $answer
				doneAfter: $doneAfter
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

func Question(question model.MediaQuizInput) (string, map[string]interface{}) {
	variables, _ := model.StructToMap(question)

	query := `
	query MediaQuiz($set: String!, $segment: String!, $theme: String!, $doneAfter: Int!, $resetProgress: Boolean!, $archiveProgress: Boolean!) {
		mediaQuiz(
			input: {
				set: $set
				segment: $segment
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
				imageUrl
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
