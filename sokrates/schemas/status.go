package schemas

import "github.com/graphql-go/graphql"

var status = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Status",
	Description: "The way to check whether backend apis are available",
	Fields: graphql.Fields{
		"overallHealth": &graphql.Field{
			Type: graphql.Boolean,
		},
		"herodotos": &graphql.Field{
			Type: health,
		},
		"sokrates": &graphql.Field{
			Type: health,
		},
		"alexandros": &graphql.Field{
			Type: health,
		},
		"dionysios": &graphql.Field{
			Type: health,
		},
	},
})

var health = graphql.NewObject(graphql.ObjectConfig{
	Name: "Health",
	Fields: graphql.Fields{
		"healthy": &graphql.Field{
			Type: graphql.Boolean,
		},
		"time": &graphql.Field{
			Type: graphql.String,
		},
		"database": &graphql.Field{
			Type: database,
		},
	},
})

var database = graphql.NewObject(graphql.ObjectConfig{
	Name: "Database",
	Fields: graphql.Fields{
		"healthy": &graphql.Field{
			Type: graphql.Boolean,
		},
		"serverName": &graphql.Field{
			Type: graphql.String,
		},
		"serverVersion": &graphql.Field{
			Type: graphql.String,
		},
		"clusterName": &graphql.Field{
			Type: graphql.String,
		},
	},
})
