package meletos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"net/http"
	"time"
)

const (
	requestID = "123456789"
	sessionID = "9999999999"
)

// ForwardGraphql sends a GraphQL query and returns a structured response
func (m *MeletosFixture) ForwardGraphql(query string, variables map[string]interface{}) (*http.Response, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	// Create a client with appropriate timeouts
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Create the request
	req, err := http.NewRequest("POST", m.sokrates.graphql, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set(config.HeaderKey, requestID)
	req.Header.Set(config.SessionIdKey, sessionID)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	return client.Do(req)
}
