package theoria

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"google.golang.org/grpc/metadata"
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
func DatabaseSpan(query map[string]interface{}, hits, timeTook int64, ctx context.Context) {
	traceID, spanID, traceCall := extractRequestIds(ctx)

	if !traceCall {
		return
	}

	parsedQuery, _ := json.Marshal(query)

	dataBaseSpan := &pb.ParabasisRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		SpanId:       spanID,
		RequestType: &pb.ParabasisRequest_DatabaseSpan{DatabaseSpan: &pb.DatabaseSpanRequest{
			Action:   "search",
			Query:    string(parsedQuery),
			Hits:     hits,
			TimeTook: timeTook,
		}},
	}

	err := streamer.Send(dataBaseSpan)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}
}

func CacheSpan(response string, sessionId string, ctx context.Context) {
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

func ServiceToServiceSpan(span *pb.ParabasisRequest, ctx context.Context) {
	traceID, spanID, traceCall := extractRequestIds(ctx)

	if !traceCall {
		return
	}

	span.TraceId = traceID
	span.ParentSpanId = spanID

	err := streamer.Send(span)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}

}
