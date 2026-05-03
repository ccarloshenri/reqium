package workflowtest

import (
	"context"
	"path/filepath"
	"testing"

	"reqium/internal/implementations/storage"
	"reqium/internal/implementations/variables"
	"reqium/internal/models"
)

func TestTemplateVariableResolverResolvesRequest(t *testing.T) {
	resolver := variables.NewTemplateVariableResolver()

	req, err := resolver.ResolveRequest(models.Request{
		Method:  "GET",
		URL:     "{{base_url}}/users",
		Headers: map[string]string{"Authorization": "Bearer {{token}}"},
		Body:    []byte(`{"account":"{{account_id}}"}`),
	}, map[string]string{
		"base_url":   "https://api.example.com",
		"token":      "abc123",
		"account_id": "acc-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.URL != "https://api.example.com/users" {
		t.Fatalf("unexpected url: %s", req.URL)
	}
	if req.Headers["Authorization"] != "Bearer abc123" {
		t.Fatalf("unexpected authorization header: %s", req.Headers["Authorization"])
	}
	if string(req.Body) != `{"account":"acc-1"}` {
		t.Fatalf("unexpected body: %s", req.Body)
	}
}

func TestJSONStorePersistsEnvironmentCollectionAndHistory(t *testing.T) {
	ctx := context.Background()
	store := storage.NewJSONStore(filepath.Join(t.TempDir(), "store.json"))
	envs := storage.NewEnvironmentRepository(store)
	collections := storage.NewCollectionRepository(store)
	history := storage.NewHistoryRepository(store)

	if err := envs.Save(ctx, models.Environment{Name: "dev", Variables: map[string]string{"base_url": "https://api.example.com"}}); err != nil {
		t.Fatalf("save env: %v", err)
	}
	if err := envs.SetActive(ctx, "dev"); err != nil {
		t.Fatalf("set active env: %v", err)
	}
	active, err := envs.Active(ctx)
	if err != nil {
		t.Fatalf("active env: %v", err)
	}
	if active.Name != "dev" {
		t.Fatalf("expected active env dev, got %s", active.Name)
	}

	if err := collections.Save(ctx, models.Collection{
		Name: "users",
		Requests: []models.SavedRequest{{
			ID:     "1",
			Name:   "list-users",
			Method: "GET",
			URL:    "{{base_url}}/users",
		}},
	}); err != nil {
		t.Fatalf("save collection: %v", err)
	}
	collection, err := collections.Get(ctx, "users")
	if err != nil {
		t.Fatalf("get collection: %v", err)
	}
	if len(collection.Requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(collection.Requests))
	}

	if err := history.Save(ctx, models.HistoryEntry{ID: "hist-1", Method: "GET", URL: "https://api.example.com/users"}); err != nil {
		t.Fatalf("save history: %v", err)
	}
	entries, err := history.List(ctx, 10)
	if err != nil {
		t.Fatalf("list history: %v", err)
	}
	if len(entries) != 1 || entries[0].ID != "hist-1" {
		t.Fatalf("unexpected history entries: %+v", entries)
	}
}
