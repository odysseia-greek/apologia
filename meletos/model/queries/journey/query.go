package journey

import "github.com/odysseia-greek/apologia/meletos/model"

func Options() string {
	return `
		query {
		  journeyOptions {
			themes {
					name
					segments {
						name
						location
						number
						coordinates{
							x
							y
						}
					}
				}
		 	 }
		}
	`
}

func Question(input model.JourneyQuizInput) (string, map[string]interface{}) {
	variables, _ := model.StructToMap(input)

	query := `
	query journeyQuiz($theme: String!, $segment: String!) {
		journeyQuiz(input: {
			theme: $theme
			segment: $segment
		}) {
			segment
			theme
			number
			translation
			sentence
			contextNote
			intro {
				author
				background
				work
			}
			quiz {
				... on MatchQuiz {
					__typename
					instruction
					pairs {
						greek
						answer
					}
				}
				... on TriviaQuiz {
					__typename
					question
					options
					answer
					note
				}
				... on StructureQuiz {
					__typename
					title
					text
					question
					options
					answer
					note
				}
				... on MediaQuiz {
					__typename
					instruction
					mediaFiles {
						word
						answer
					}
				}
				... on FinalTranslationQuiz {
					__typename
					instruction
					options
					answer
				}
			}
		}
	}
	`

	return query, variables
}
