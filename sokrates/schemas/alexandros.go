package schemas

import "github.com/graphql-go/graphql"

var dictionary = graphql.NewObject(graphql.ObjectConfig{
	Name: "Hit",
	Fields: graphql.Fields{
		"greek": &graphql.Field{
			Type: graphql.String,
		},
		"english": &graphql.Field{
			Type: graphql.String,
		},
		"dutch": &graphql.Field{
			Type: graphql.String,
		},
		"linkedWord": &graphql.Field{
			Type: graphql.String,
		},
		"original": &graphql.Field{
			Type: graphql.String,
		},
	},
})
