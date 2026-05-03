package requestinputtest

import (
	"errors"
	"reflect"
	"testing"

	"reqium/internal/requestinput"
)

type fakeFileReader struct {
	data []byte
	err  error
}

func (r fakeFileReader) Read(path string) ([]byte, error) {
	return r.data, r.err
}

func TestParseHeaders(t *testing.T) {
	headers, err := requestinput.ParseHeaders([]string{
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
	_, err := requestinput.ParseHeaders([]string{"Authorization"})
	if err == nil {
		t.Fatal("expected invalid header error")
	}
}

func TestLoadBodyFromRawString(t *testing.T) {
	body, err := requestinput.LoadBody(requestinput.BodyOptions{Body: `{"name":"John"}`}, fakeFileReader{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != `{"name":"John"}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestLoadBodyFromFile(t *testing.T) {
	reader := fakeFileReader{data: []byte(`{"name":"Jane"}`)}
	body, err := requestinput.LoadBody(requestinput.BodyOptions{BodyFile: "payload.json"}, reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != `{"name":"Jane"}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestLoadBodyRejectsRawAndFileTogether(t *testing.T) {
	_, err := requestinput.LoadBody(requestinput.BodyOptions{Body: "{}", BodyFile: "payload.json"}, fakeFileReader{})
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestBuildRequestRejectsBodyForGet(t *testing.T) {
	_, err := requestinput.BuildRequest("GET", "https://api.example.com/users", nil, requestinput.BodyOptions{
		Body: "{}",
	}, 30, fakeFileReader{})
	if err == nil {
		t.Fatal("expected body method error")
	}
}

func TestLoadBodyReturnsFileError(t *testing.T) {
	expected := errors.New("read failed")
	_, err := requestinput.LoadBody(requestinput.BodyOptions{BodyFile: "payload.json"}, fakeFileReader{err: expected})
	if !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
}
