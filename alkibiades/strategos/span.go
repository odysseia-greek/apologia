package strategos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc/metadata"
	"strings"
)

func extractRequestIds(ctx context.Context) (string, string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	var requestId string
	if ok {
		headerValue := md.Get(service.HeaderKey)
		if len(headerValue) > 0 {
			requestId = headerValue[0]
		}
	}

	splitID := strings.Split(requestId, "+")

	traceCall := false
	var traceID, spanID string

	if len(splitID) >= 3 {
		traceCall = splitID[2] == "1"
	}

	if len(splitID) >= 1 {
		traceID = splitID[0]
	}
	if len(splitID) >= 2 {
		spanID = splitID[1]
	}

	return traceID, spanID, traceCall
}
func databaseSpan(response *models.Response, query map[string]interface{}, ctx context.Context) {
	traceID, spanID, traceCall := extractRequestIds(ctx)

	if response == nil || !traceCall {
		return
	}

	parsedQuery, _ := json.Marshal(query)
	hits := int64(0)
	if response.Hits.Hits != nil {
		hits = response.Hits.Total.Value
	}

	dataBaseSpan := &pb.ParabasisRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		SpanId:       spanID,
		RequestType: &pb.ParabasisRequest_DatabaseSpan{DatabaseSpan: &pb.DatabaseSpanRequest{
			Action:   "search",
			Query:    string(parsedQuery),
			Hits:     hits,
			TimeTook: response.Took,
		}},
	}

	err := streamer.Send(dataBaseSpan)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}
}

func cacheSpan(response string, sessionId string, ctx context.Context) {
	traceID, spanID, traceCall := extractRequestIds(ctx)

	if !traceCall {
		return
	}

	span := &pb.ParabasisRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		SpanId:       spanID,
		RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
			Action: fmt.Sprintf("taken from cache with key: %s", sessionId),
			Status: response,
		}},
	}

	err := streamer.Send(span)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}
}
