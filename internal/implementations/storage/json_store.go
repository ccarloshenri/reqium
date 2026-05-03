package storage

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	reqerrors "reqium/internal/errors"
	"reqium/internal/models"
)

type database struct {
	History      []models.HistoryEntry `json:"history"`
	Environments []models.Environment  `json:"environments"`
	Collections  []models.Collection   `json:"collections"`
}

type JSONStore struct {
	path string
	mu   sync.Mutex
}

func NewDefaultJSONStore() (*JSONStore, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(base, "reqium", "store.json")
	return NewJSONStore(path), nil
}

func NewJSONStore(path string) *JSONStore {
	return &JSONStore{path: path}
}

func (s *JSONStore) SaveHistory(ctx context.Context, entry models.HistoryEntry) error {
	return s.update(ctx, func(db *database) error {
		db.History = append([]models.HistoryEntry{entry}, db.History...)
		return nil
	})
}

func (s *JSONStore) ListHistory(ctx context.Context, limit int) ([]models.HistoryEntry, error) {
	db, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > len(db.History) {
		limit = len(db.History)
	}
	return slices.Clone(db.History[:limit]), nil
}

func (s *JSONStore) GetHistory(ctx context.Context, id string) (models.HistoryEntry, error) {
	db, err := s.load(ctx)
	if err != nil {
		return models.HistoryEntry{}, err
	}
	for _, entry := range db.History {
		if entry.ID == id {
			return entry, nil
		}
	}
	return models.HistoryEntry{}, reqerrors.ErrNotFound
}

func (s *JSONStore) SaveEnvironment(ctx context.Context, env models.Environment) error {
	return s.update(ctx, func(db *database) error {
		now := time.Now()
		for i, existing := range db.Environments {
			if existing.Name == env.Name {
				env.CreatedAt = existing.CreatedAt
				env.UpdatedAt = now
				db.Environments[i] = env
				return nil
			}
		}
		env.CreatedAt = now
		env.UpdatedAt = now
		db.Environments = append(db.Environments, env)
		return nil
	})
}

func (s *JSONStore) ListEnvironments(ctx context.Context) ([]models.Environment, error) {
	db, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	return slices.Clone(db.Environments), nil
}

func (s *JSONStore) GetEnvironment(ctx context.Context, name string) (models.Environment, error) {
	db, err := s.load(ctx)
	if err != nil {
		return models.Environment{}, err
	}
	for _, env := range db.Environments {
		if env.Name == name {
			return env, nil
		}
	}
	return models.Environment{}, reqerrors.ErrNotFound
}

func (s *JSONStore) DeleteEnvironment(ctx context.Context, name string) error {
	return s.update(ctx, func(db *database) error {
		for i, env := range db.Environments {
			if env.Name == name {
				db.Environments = append(db.Environments[:i], db.Environments[i+1:]...)
				return nil
			}
		}
		return reqerrors.ErrNotFound
	})
}

func (s *JSONStore) SetActiveEnvironment(ctx context.Context, name string) error {
	return s.update(ctx, func(db *database) error {
		found := false
		for i := range db.Environments {
			active := db.Environments[i].Name == name
			db.Environments[i].Active = active
			if active {
				found = true
				db.Environments[i].UpdatedAt = time.Now()
			}
		}
		if !found {
			return reqerrors.ErrNotFound
		}
		return nil
	})
}

func (s *JSONStore) ActiveEnvironment(ctx context.Context) (models.Environment, error) {
	db, err := s.load(ctx)
	if err != nil {
		return models.Environment{}, err
	}
	for _, env := range db.Environments {
		if env.Active {
			return env, nil
		}
	}
	return models.Environment{}, reqerrors.ErrNotFound
}

func (s *JSONStore) SaveCollection(ctx context.Context, collection models.Collection) error {
	return s.update(ctx, func(db *database) error {
		now := time.Now()
		for i, existing := range db.Collections {
			if existing.Name == collection.Name {
				collection.CreatedAt = existing.CreatedAt
				collection.UpdatedAt = now
				db.Collections[i] = collection
				return nil
			}
		}
		collection.CreatedAt = now
		collection.UpdatedAt = now
		db.Collections = append(db.Collections, collection)
		return nil
	})
}

func (s *JSONStore) ListCollections(ctx context.Context) ([]models.Collection, error) {
	db, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	return slices.Clone(db.Collections), nil
}

func (s *JSONStore) GetCollection(ctx context.Context, name string) (models.Collection, error) {
	db, err := s.load(ctx)
	if err != nil {
		return models.Collection{}, err
	}
	for _, collection := range db.Collections {
		if collection.Name == name {
			return collection, nil
		}
	}
	return models.Collection{}, reqerrors.ErrNotFound
}

func (s *JSONStore) DeleteCollection(ctx context.Context, name string) error {
	return s.update(ctx, func(db *database) error {
		for i, collection := range db.Collections {
			if collection.Name == name {
				db.Collections = append(db.Collections[:i], db.Collections[i+1:]...)
				return nil
			}
		}
		return reqerrors.ErrNotFound
	})
}

func (s *JSONStore) load(ctx context.Context) (database, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadLocked(ctx)
}

func (s *JSONStore) update(ctx context.Context, fn func(*database) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	db, err := s.loadLocked(ctx)
	if err != nil {
		return err
	}
	if err := fn(&db); err != nil {
		return err
	}
	return s.saveLocked(ctx, db)
}

func (s *JSONStore) loadLocked(ctx context.Context) (database, error) {
	if err := ctx.Err(); err != nil {
		return database{}, err
	}

	body, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return database{}, nil
		}
		return database{}, err
	}
	if len(body) == 0 {
		return database{}, nil
	}

	var db database
	if err := json.Unmarshal(body, &db); err != nil {
		return database{}, err
	}
	return db, nil
}

func (s *JSONStore) saveLocked(ctx context.Context, db database) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	body, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, body, 0o600)
}
