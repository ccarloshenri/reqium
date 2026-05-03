package storage

import (
	"context"

	"reqium/internal/models"
)

type HistoryRepository struct {
	store *JSONStore
}

func NewHistoryRepository(store *JSONStore) *HistoryRepository {
	return &HistoryRepository{store: store}
}

func (r *HistoryRepository) Save(ctx context.Context, entry models.HistoryEntry) error {
	return r.store.SaveHistory(ctx, entry)
}

func (r *HistoryRepository) List(ctx context.Context, limit int) ([]models.HistoryEntry, error) {
	return r.store.ListHistory(ctx, limit)
}

func (r *HistoryRepository) Get(ctx context.Context, id string) (models.HistoryEntry, error) {
	return r.store.GetHistory(ctx, id)
}

type EnvironmentRepository struct {
	store *JSONStore
}

func NewEnvironmentRepository(store *JSONStore) *EnvironmentRepository {
	return &EnvironmentRepository{store: store}
}

func (r *EnvironmentRepository) Save(ctx context.Context, env models.Environment) error {
	return r.store.SaveEnvironment(ctx, env)
}

func (r *EnvironmentRepository) List(ctx context.Context) ([]models.Environment, error) {
	return r.store.ListEnvironments(ctx)
}

func (r *EnvironmentRepository) Get(ctx context.Context, name string) (models.Environment, error) {
	return r.store.GetEnvironment(ctx, name)
}

func (r *EnvironmentRepository) Delete(ctx context.Context, name string) error {
	return r.store.DeleteEnvironment(ctx, name)
}

func (r *EnvironmentRepository) SetActive(ctx context.Context, name string) error {
	return r.store.SetActiveEnvironment(ctx, name)
}

func (r *EnvironmentRepository) Active(ctx context.Context) (models.Environment, error) {
	return r.store.ActiveEnvironment(ctx)
}

type CollectionRepository struct {
	store *JSONStore
}

func NewCollectionRepository(store *JSONStore) *CollectionRepository {
	return &CollectionRepository{store: store}
}

func (r *CollectionRepository) Save(ctx context.Context, collection models.Collection) error {
	return r.store.SaveCollection(ctx, collection)
}

func (r *CollectionRepository) List(ctx context.Context) ([]models.Collection, error) {
	return r.store.ListCollections(ctx)
}

func (r *CollectionRepository) Get(ctx context.Context, name string) (models.Collection, error) {
	return r.store.GetCollection(ctx, name)
}

func (r *CollectionRepository) Delete(ctx context.Context, name string) error {
	return r.store.DeleteCollection(ctx, name)
}
