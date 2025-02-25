package schemas

import "github.com/graphql-go/graphql"

var analyzeTextResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AnalyzeTextResponse",
	Fields: graphql.Fields{
		"rootword": &graphql.Field{Type: graphql.String},
		"conjugations": &graphql.Field{
			Type: graphql.NewList(conjugationResponseType),
		},
		"results": &graphql.Field{
			Type: graphql.NewList(analyzeResultType),
		},
	},
})

var conjugationResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ConjugationResponse",
	Fields: graphql.Fields{
		"word": &graphql.Field{Type: graphql.String},
		"rule": &graphql.Field{Type: graphql.String},
	},
})

// Define AnalyzeResult Type
var analyzeResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AnalyzeResult",
	Fields: graphql.Fields{
		"referenceLink": &graphql.Field{Type: graphql.String},
		"text":          &graphql.Field{Type: rhemaType},
		"author":        &graphql.Field{Type: graphql.String},
		"book":          &graphql.Field{Type: graphql.String},
		"reference":     &graphql.Field{Type: graphql.String},
	},
})

// Define Rhema Type
var rhemaType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Rhema",
	Fields: graphql.Fields{
		"greek": &graphql.Field{Type: graphql.String},
		"translations": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"section": &graphql.Field{Type: graphql.String},
	},
})
