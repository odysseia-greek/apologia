package dialogue

import "github.com/odysseia-greek/apologia/meletos/model"

func Options() string {
	return `
		query {
		  dialogueOptions {
			themes {
					name
					maxSet
				}
		  }
		}
	`
}

// Question returns a GraphQL query for dialogue quiz and the variables needed for the query
func Question(input model.DialogueQuizInput) (string, map[string]interface{}) {
	variables, _ := model.StructToMap(input)

	query := `
	query dialogueQuiz($theme: String!, $set: String!) {
		dialogueQuiz(input: {
			theme: $theme
			set: $set
		}) {
			quizMetadata {
				language
			}
			segment
			theme
			set
			reference
			dialogue {
				introduction
				linkToPerseus
				section
				speakers {
					name
					shorthand
					translation
				}
			}
			content {
				translation
				greek
				place
				speaker
			}
		}
	}
	`

	return query, variables
}

// Answer returns a GraphQL query for dialogue answer and the variables needed for the query
func Answer(input model.DialogueAnswerInput) (string, map[string]interface{}) {
	variables, _ := model.StructToMap(input)

	query := `
	query dialogueAnswer($theme: String!, $set: String!, $content: [DialogueInputContent!]!) {
		dialogueAnswer(input: {
			theme: $theme
			set: $set
			content: $content
		}) {
			percentage
			input {
				greek
				place
				speaker
				translation
			}
			answer {
				greek
				place
				speaker
				translation
			}
			wronglyPlaced {
				greek
				translation
				speaker
				place
				correctPlace
			}
		}
	}
	`

	return query, variables
}
