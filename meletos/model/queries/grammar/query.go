package grammar

import (
	"github.com/odysseia-greek/apologia/meletos/model"
)

func Options() string {
	return `
		query {
		  grammarOptions {
			themes {
					name
					segments {
						name
						difficulty
						maxSet
					}
				}
		  }
		}
	`
}

func Answer(answer model.GrammarAnswerInput) (string, map[string]interface{}) {
	variables, _ := model.StructToMap(answer)

	query := `
	query grammarAnswer($set: String!, $segment: String!, $theme: String!, $quizWord: String!, $answer: String!, $doneAfter: Int!, $comprehensive: Boolean!, $dictionaryForm: String!) {
		grammarAnswer(
			input: {
				set: $set
				segment: $segment
				theme: $theme
				quizWord: $quizWord	
				answer: $answer
				doneAfter: $doneAfter
				comprehensive: $comprehensive
				dictionaryForm: $dictionaryForm
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

func Question(question model.GrammarQuizInput) (string, map[string]interface{}) {
	variables, _ := model.StructToMap(question)

	query := `
	query GrammarQuiz($set: String!, $segment: String!, $theme: String!, $doneAfter: Int!, $resetProgress: Boolean!, $archiveProgress: Boolean!) {
		grammarQuiz(
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
				stem
				difficulty
				description
				contractionRule
				dictionaryForm
				translation
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
