package formatter

import (
	"strings"
	"testing"

	"reqium/internal/domain"
)

func TestJSONFormatterPrettyFormatsValidJSON(t *testing.T) {
	body, ok := NewJSONFormatter().Pretty([]byte(`{"name":"John","roles":["admin"]}`))
	if !ok {
		t.Fatal("expected valid json to be formatted")
	}

	expected := "{\n  \"name\": \"John\",\n  \"roles\": [\n    \"admin\"\n  ]\n}"
	if string(body) != expected {
		t.Fatalf("expected %q, got %q", expected, body)
	}
}

func TestJSONFormatterLeavesInvalidJSON(t *testing.T) {
	input := []byte(`{"name":`)
	body, ok := NewJSONFormatter().Pretty(input)
	if ok {
		t.Fatal("expected invalid json")
	}
	if string(body) != string(input) {
		t.Fatalf("expected original body, got %q", body)
	}
}

func TestResponseFormatterIncludesStatusHeadersBodyAndDuration(t *testing.T) {
	output, err := NewResponseFormatter(true).Format(domain.Response{
		StatusCode: 201,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: []byte(`{"ok":true}`),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, value := range []string{"Status: 201", "Headers:", "Content-Type: application/json", "\"ok\": true", "Duration:"} {
		if !strings.Contains(output, value) {
			t.Fatalf("expected output to contain %q, got %s", value, output)
		}
	}
}
