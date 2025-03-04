package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"io"
	"net/http"
	"strings"

	"github.com/odysseia-greek/agora/plato/logging"
)

type Adapter func(http.Handler) http.Handler

// Adapt Iterate over adapters and run them one by one
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func SetCorsHeaders() Adapter {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			// Allow any origin in the allowedOrigins slice
			allowedOrigins := []string{"localhost"}

			for _, o := range allowedOrigins {
				if strings.Contains(origin, o) {
					logging.Debug(fmt.Sprintf("setting CORS header for origin: %s", origin))
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					if r.Method == "OPTIONS" {
						w.WriteHeader(http.StatusOK)
						return
					}
					break
				}
			}
			f.ServeHTTP(w, r)
		})
	}
}

// LogRequestDetails is a middleware function that captures and logs details of incoming requests,
// and initiates traces based on the configured trace probabilities for specific GraphQL operations.
// It reads the incoming request body to extract the operation name and query from GraphQL requests.
// The middleware then checks the trace configuration to determine whether to initiate a trace for
// the given operation. If the trace probability condition is met, a trace is started using the
// provided tracer's StartTrace method. The trace ID is logged, and the middleware creates a new
// context with the trace ID to pass it along to downstream handlers.
//
// Parameters:
// - tracer: The tracer instance used to initiate traces.
// - traceConfig: The configuration specifying the trace probabilities for specific operations.
//
// Returns:
// An Adapter that wraps an http.Handler and performs the described middleware actions.
func LogRequestDetails(tracer pb.TraceService_ChorusClient) Adapter {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(config.HeaderKey)
			trace := traceFromString(requestId)

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // Set the original request body

			var bodyClone map[string]interface{}
			decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
			if err := decoder.Decode(&bodyClone); err != nil {
				http.Error(w, "Failed to parse request body", http.StatusInternalServerError)
				return
			}

			operationName, _ := bodyClone["operationName"].(string)
			query, _ := bodyClone["query"].(string)

			if operationName == "" {
				splitQuery := strings.Split(query, "{")
				if len(splitQuery) != 0 {
					if strings.Contains(splitQuery[0], " ") {
						operationName = strings.Split(splitQuery[0], " ")[1]
						logging.Debug(fmt.Sprintf("extracted operationName from query: %s", operationName))
					}
				}
			}

			spanID := comedy.GenerateSpanID()

			payload := &pb.TraceRequestWithBody{
				Method:    r.Method,
				Url:       r.URL.RequestURI(),
				Host:      r.Host,
				Operation: operationName,
				RootQuery: query,
			}

			parabasis := &pb.ParabasisRequest{
				TraceId:      trace.TraceId,
				ParentSpanId: trace.SpanId,
				SpanId:       spanID,
				RequestType: &pb.ParabasisRequest_TraceBody{
					TraceBody: payload,
				},
			}
			if err := tracer.Send(parabasis); err != nil {
				logging.Error(fmt.Sprintf("failed to send trace data: %v", err))
			}

			logging.Trace(fmt.Sprintf("trace with requestID: %s and parentSpan: %s and span: %s", trace.TraceId, trace.SpanId, spanID))

			ctx := context.WithValue(r.Context(), config.HeaderKey, requestId)
			f.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func traceFromString(requestId string) *pb.TraceBare {
	splitID := strings.Split(requestId, "+")

	trace := &pb.TraceBare{}

	if len(splitID) >= 3 {
		trace.Save = splitID[2] == "1"
	}

	if len(splitID) >= 1 {
		trace.TraceId = splitID[0]
	}
	if len(splitID) >= 2 {
		trace.SpanId = splitID[1]
	}

	return trace
}
