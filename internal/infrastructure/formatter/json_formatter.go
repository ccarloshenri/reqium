package formatter

import (
	"bytes"
	"encoding/json"
)

type JSONFormatter struct{}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

func (f *JSONFormatter) Pretty(body []byte) ([]byte, bool) {
	if !json.Valid(body) {
		return body, false
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, body, "", "  "); err != nil {
		return body, false
	}

	return pretty.Bytes(), true
}
