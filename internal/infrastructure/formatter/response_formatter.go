package formatter

import (
	"fmt"
	"slices"
	"strings"

	"reqium/internal/domain"
)

type ResponseFormatter struct {
	prettyJSON bool
	json       *JSONFormatter
}

func NewResponseFormatter(prettyJSON bool) *ResponseFormatter {
	return &ResponseFormatter{
		prettyJSON: prettyJSON,
		json:       NewJSONFormatter(),
	}
}

func (f *ResponseFormatter) Format(response domain.Response) (string, error) {
	var builder strings.Builder

	fmt.Fprintf(&builder, "Status: %d\n", response.StatusCode)
	fmt.Fprintf(&builder, "Duration: %s\n", response.Duration)
	builder.WriteString("Headers:\n")

	keys := make([]string, 0, len(response.Headers))
	for key := range response.Headers {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	for _, key := range keys {
		for _, value := range response.Headers[key] {
			fmt.Fprintf(&builder, "  %s: %s\n", key, value)
		}
	}

	body := response.Body
	if f.prettyJSON {
		body, _ = f.json.Pretty(response.Body)
	}

	builder.WriteString("Body:\n")
	builder.Write(body)
	if len(body) > 0 && body[len(body)-1] != '\n' {
		builder.WriteByte('\n')
	}

	return builder.String(), nil
}
