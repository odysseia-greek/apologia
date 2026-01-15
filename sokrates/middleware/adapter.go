package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
)

type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// statusRecorder lets us capture the final HTTP status code written by the handler.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	// If handler never called WriteHeader, status is implicitly 200 on first Write.
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}

func LogRequestDetails(tracer arv1.TraceService_ChorusClient) Adapter {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(config.HeaderKey)
			sessionId := r.Header.Get(config.SessionIdKey)

			trace := comedy.TraceBareFromString(requestId)
			// If this request isn't being traced, just pass through.
			if trace.TraceId == "" || trace.SpanId == "" || !trace.Save {
				f.ServeHTTP(w, r)
				return
			}

			// Read body once; restore for downstream
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			_ = r.Body.Close()
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			var (
				operationName string
				query         string
			)

			// Best-effort JSON parse (GraphQL POST)
			if len(bodyBytes) > 0 {
				var bodyClone map[string]any
				if err := json.Unmarshal(bodyBytes, &bodyClone); err == nil {
					if v, ok := bodyClone["operationName"].(string); ok {
						operationName = v
					}
					if v, ok := bodyClone["query"].(string); ok {
						query = v
					}
				}
			}

			// Fallback extraction if operationName missing but query exists
			if operationName == "" && query != "" {
				// naive, but works for "query foo(...) {"
				splitQuery := strings.Split(query, "{")
				if len(splitQuery) > 1 {
					part := splitQuery[1]
					if strings.Contains(part, "(") {
						part = strings.Split(part, "(")[0]
					}
					operationName = strings.TrimSpace(part)
				}
			}

			// Create a span representing THIS GraphQL operation
			parentSpan := trace.SpanId
			graphqlSpan := comedy.GenerateSpanID()

			payload := &arv1.ObserveGraphQL{
				Operation: operationName,
				RootQuery: query,
			}

			// Emit GRAPHQL event (child of incoming span)
			parabasis := &arv1.ObserveRequest{
				TraceId:      trace.TraceId,
				ParentSpanId: parentSpan,
				SpanId:       graphqlSpan,
				Kind:         &arv1.ObserveRequest_Graphql{Graphql: payload},
			}

			if err := tracer.Send(parabasis); err != nil {
				logging.Error(fmt.Sprintf("failed to send graphql trace data: %v", err))
			}

			if operationName != "IntrospectionQuery" {
				if jsonPayload, err := json.MarshalIndent(payload, "", "  "); err == nil {
					logging.Info(fmt.Sprintf("REQUEST | traceId: %s and params:\n%s", trace.TraceId, string(jsonPayload)))
				}
			}

			// Propagate: downstream should treat this GraphQL span as the current parent
			trace.SpanId = graphqlSpan
			newRequestId := comedy.CreateCombinedId(trace)

			// You usually don't want to mutate response headers here (clients may ignore),
			// but keeping your behavior for now:
			w.Header().Set(config.HeaderKey, newRequestId)
			w.Header().Set(config.SessionIdKey, sessionId)

			ctx := context.WithValue(r.Context(), config.HeaderKey, newRequestId)
			ctx = context.WithValue(ctx, config.SessionIdKey, sessionId)

			// Wrap writer to capture status + measure duration
			rec := &statusRecorder{ResponseWriter: w}
			start := time.Now()

			f.ServeHTTP(rec, r.WithContext(ctx))

			dur := time.Since(start)
			status := rec.status
			if status == 0 {
				status = http.StatusOK
			}

			// Emit close-hop (TRACE_HOP_STOP) for the GraphQL span
			// ParentSpanId should be the span we are closing (graphqlSpan),
			// and SpanId should be a fresh span id for the stop event itself.
			stop := &arv1.ObserveRequest{
				TraceId:      trace.TraceId,
				ParentSpanId: graphqlSpan,
				SpanId:       comedy.GenerateSpanID(),
				Kind: &arv1.ObserveRequest_TraceHopStop{
					TraceHopStop: &arv1.ObserveTraceHopStop{
						ResponseCode: int32(status),      // HTTP code
						TookMs:       dur.Milliseconds(), // duration in ms
					},
				},
			}

			if err := tracer.Send(stop); err != nil {
				logging.Error(fmt.Sprintf("failed to send graphql close-hop: %v", err))
			}

			logging.Trace(fmt.Sprintf(
				"graphql span closed | traceId=%s parent=%s span=%s status=%d tookMs=%d",
				trace.TraceId, parentSpan, graphqlSpan, status, dur.Milliseconds(),
			))
		})
	}
}
