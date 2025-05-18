package routing

import (
	"encoding/json"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/mux"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/apologia/sokrates/gateway"
	"github.com/odysseia-greek/apologia/sokrates/graph"
	"github.com/odysseia-greek/apologia/sokrates/middleware"
	"net/http"
	"os"
	"time"
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
	srv.Use(extension.Introspection{})

	graphqlHandler := middleware.Adapt(
		srv,
		middleware.LogRequestDetails(handlerConfig.Streamer),
	)

	serveMux.Handle("/sokrates/graphql", graphqlHandler)

	// --- health endpoints ---
	serveMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeHealthResponse(w)
	})
	serveMux.HandleFunc("/sokrates/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		writeHealthResponse(w)
	})

	return serveMux
}

// writeHealthResponse is the lightweight ping response
func writeHealthResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	resp := models.Health{
		Healthy: true,
		Time:    time.Now().Format(time.RFC3339),
		Version: os.Getenv("VERSION"),
	}
	json.NewEncoder(w).Encode(resp)
}
