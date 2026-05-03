package cli

import (
	"errors"
	"reflect"
	"testing"
)

type fakeFileReader struct {
	data []byte
	err  error
}

func (r fakeFileReader) Read(path string) ([]byte, error) {
	return r.data, r.err
}

func TestParseHeaders(t *testing.T) {
	headers, err := parseHeaders([]string{
		"Authorization: Bearer token",
		"Content-Type: application/json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]string{
		"Authorization": "Bearer token",
		"Content-Type":  "application/json",
	}
	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("expected %v, got %v", expected, headers)
	}
}

func TestParseHeadersRejectsInvalidHeader(t *testing.T) {
	_, err := parseHeaders([]string{"Authorization"})
	if err == nil {
		t.Fatal("expected invalid header error")
	}
}

func TestLoadBodyFromRawString(t *testing.T) {
	body, err := loadBody(requestOptions{body: `{"name":"John"}`}, fakeFileReader{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != `{"name":"John"}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestLoadBodyFromFile(t *testing.T) {
	reader := fakeFileReader{data: []byte(`{"name":"Jane"}`)}
	body, err := loadBody(requestOptions{bodyFile: "payload.json"}, reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != `{"name":"Jane"}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestLoadBodyRejectsRawAndFileTogether(t *testing.T) {
	_, err := loadBody(requestOptions{body: "{}", bodyFile: "payload.json"}, fakeFileReader{})
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestBuildRequestRejectsBodyForGet(t *testing.T) {
	_, err := buildRequest("GET", "https://api.example.com/users", requestOptions{
		body:       "{}",
		timeoutSec: 30,
	}, fakeFileReader{})
	if err == nil {
		t.Fatal("expected body method error")
	}
}

func TestLoadBodyReturnsFileError(t *testing.T) {
	expected := errors.New("read failed")
	_, err := loadBody(requestOptions{bodyFile: "payload.json"}, fakeFileReader{err: expected})
	if !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
}
