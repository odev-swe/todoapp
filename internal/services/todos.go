package services

import (
	"context"

	"github.com/odev-swe/todoapp/internal/store"
	"github.com/odev-swe/todoapp/internal/types"
)

type TodosService struct {
	store *store.TodosStore
}

func NewTodosService(store *store.TodosStore) *TodosService {
	return &TodosService{store: store}
}

func (s *TodosService) Get(ctx context.Context) ([]types.Todos, error) {
	return s.store.Get(ctx)
}

func (s *TodosService) Create(ctx context.Context, req types.TodosPostRequestBody) (*types.Todos, error) {
	return s.store.Create(ctx, req)
}

func (s *TodosService) Update(ctx context.Context, req types.TodosPutRequestBody) (*types.Todos, error) {
	return s.store.Update(ctx, req)
}

func (s *TodosService) Delete(ctx context.Context, req types.TodosDeleteRequestBody) error {
	return s.store.Delete(ctx, req)
}
