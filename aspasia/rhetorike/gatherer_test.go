package rhetorike

import (
	"errors"
	"testing"

	dionysiosv1 "github.com/odysseia-greek/alexandreia/dionysios/gen/go/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapDionysiosResearch(t *testing.T) {
	source := &dionysiosv1.ResearchResponse{
		Rootword:     "λύω",
		PartOfSpeech: "verb",
		Conjugations: []*dionysiosv1.Conjugation{{Word: "λύει", Rule: "present active indicative"}},
		Results: []*dionysiosv1.AnalyzeResult{{
			ReferenceLink: "https://example.test/text",
			Author:        "Plato",
			Book:          "Republic",
			Reference:     "1.1",
			Text: &dionysiosv1.Rhema{
				Greek:        "λύει",
				Translations: []string{"he releases"},
				Section:      "1",
			},
		}},
	}

	result := mapDionysiosResearch(source)

	if result.Rootword != source.Rootword || result.PartOfSpeech != source.PartOfSpeech {
		t.Fatalf("unexpected research metadata: %#v", result)
	}
	if len(result.Conjugations) != 1 || result.Conjugations[0].Word != "λύει" {
		t.Fatalf("unexpected conjugations: %#v", result.Conjugations)
	}
	if len(result.Texts) != 1 || result.Texts[0].Text == nil {
		t.Fatalf("unexpected text results: %#v", result.Texts)
	}
	if result.Texts[0].Text.Greek != "λύει" || result.Texts[0].Text.Translations[0] != "he releases" {
		t.Fatalf("unexpected mapped text: %#v", result.Texts[0].Text)
	}
}

func TestIsDionysiosNoResults(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "canonical not found", err: status.Error(codes.NotFound, "no research results"), want: true},
		{name: "v0.3.3 wrapped no hits", err: status.Error(codes.InvalidArgument, `research failed: no hits for rootword "ἄγνωστος"`), want: true},
		{name: "other invalid argument", err: status.Error(codes.InvalidArgument, "rootword is required"), want: false},
		{name: "service unavailable", err: status.Error(codes.Unavailable, "connection refused"), want: false},
		{name: "ordinary error", err: errors.New("no hits for rootword"), want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isDionysiosNoResults(test.err); got != test.want {
				t.Fatalf("isDionysiosNoResults() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestMapDionysiosResearchHandlesNil(t *testing.T) {
	result := mapDionysiosResearch(nil)
	if result == nil {
		t.Fatal("expected an empty response")
	}
	if result.Conjugations == nil || result.Texts == nil {
		t.Fatalf("expected initialized result slices: %#v", result)
	}
}
