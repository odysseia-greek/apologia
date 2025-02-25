package routing

import (
	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	plato "github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/apologia/sokrates/gateway"
	"github.com/odysseia-greek/apologia/sokrates/middleware"
	"github.com/odysseia-greek/apologia/sokrates/schemas"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(tracer pb.TraceService_ChorusClient) *mux.Router {
	serveMux := mux.NewRouter()

	srv := handler.New(&handler.Config{
		Schema:   &schemas.HomerosSchema,
		Pretty:   true,
		GraphiQL: false,
	})

	serveMux.HandleFunc("/homeros/v1/health", plato.Adapt(gateway.HealthProbe, plato.ValidateRestMethod("GET"), plato.SetCorsHeaders()))
	serveMux.Handle("/graphql", middleware.Adapt(srv, middleware.LogRequestDetails(tracer), middleware.SetCorsHeaders()))

	return serveMux
}
