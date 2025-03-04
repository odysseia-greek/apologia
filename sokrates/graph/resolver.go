package graph

import (
	"github.com/odysseia-greek/apologia/sokrates/gateway"
)

// Resolver struct for dependency injection
type Resolver struct {
	Handler *gateway.SokratesHandler
}
