package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Todos struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	DueDate     time.Time `json:"due_date,omitempty" `
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TodosPostRequestBody struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	DueDate     time.Time `json:"due_date,omitempty" default:"2022-01-01T00:00:00Z"`
}

type TodosPutRequestBody struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	DueDate     time.Time `json:"due_date,omitempty" default:"2022-01-01T00:00:00Z"`
}

type TodosDeleteRequestBody struct {
	Id uuid.UUID `json:"id"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type TodosServices interface {
	Create(ctx context.Context, req TodosPostRequestBody) (*Todos, error)
	Get(ctx context.Context) ([]Todos, error)
	Update(ctx context.Context, req TodosPutRequestBody) (*Todos, error)
	Delete(ctx context.Context, req TodosDeleteRequestBody) error
}
