package routing

import (
	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"github.com/odysseia-greek/apologia/sokrates/middleware"
	"github.com/odysseia-greek/apologia/sokrates/schemas"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes(tracer pb.TraceService_ChorusClient) *mux.Router {
	serveMux := mux.NewRouter()

	srv := handler.New(&handler.Config{
		Schema:   &schemas.SokratesSchema,
		Pretty:   true,
		GraphiQL: false,
	})

	serveMux.Handle("/sokrates/graphql", middleware.Adapt(srv, middleware.LogRequestDetails(tracer), middleware.SetCorsHeaders()))

	return serveMux
}
