package routing

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/apologia/sokrates/gateway"
	"github.com/odysseia-greek/apologia/sokrates/graph"
	"github.com/odysseia-greek/apologia/sokrates/middleware"
	"github.com/vektah/gqlparser/v2/ast"
)

// InitRoutes initializes the mux router with middleware and GraphQL handler
func InitRoutes(handlerConfig *gateway.SokratesHandler) *mux.Router {
	serveMux := mux.NewRouter()

	srv := handler.New(graph.NewExecutableSchema(
		graph.Config{Resolvers: &graph.Resolver{Handler: handlerConfig}},
	))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](100))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	graphqlHandler := middleware.Adapt(
		srv,
		middleware.LogRequestDetails(handlerConfig.Streamer),
	)

	serveMux.Handle("/sokrates/graphql", graphqlHandler)

	return serveMux
}
