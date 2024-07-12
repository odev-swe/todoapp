package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/odev-swe/todoapp/internal/types"
)

type TodosStore struct {
	db *pgxpool.Pool
}

func NewTodosStore(db *pgxpool.Pool) *TodosStore {
	return &TodosStore{
		db: db,
	}
}

func (s *TodosStore) Create(ctx context.Context, req types.TodosPostRequestBody) (*types.Todos, error) {
	// acquire connection
	conn, err := s.db.Acquire(ctx)

	if err != nil {
		return nil, err
	}
	// defer release connection
	defer conn.Release()

	// get user id from context
	userId := ctx.Value("user-id")
	uuidUserId, err := uuid.Parse(userId.(string))

	if err != nil {
		return nil, err
	}

	// perform query
	var todo types.Todos

	prepareQuery := "INSERT INTO todos (title, description, completed, due_date, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, description, completed, due_date, created_at, updated_at"

	err = conn.QueryRow(ctx, prepareQuery, req.Title, req.Description, req.Completed, req.DueDate, uuidUserId).Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed, &todo.DueDate, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s *TodosStore) Get(ctx context.Context) ([]types.Todos, error) {
	// acquire connection
	conn, err := s.db.Acquire(ctx)

	if err != nil {
		return nil, err
	}
	// defer release connection
	defer conn.Release()

	// get user id from context
	userId := ctx.Value("user-id")
	uuidUserId, err := uuid.Parse(userId.(string))

	limit := ctx.Value("limit")
	offset := ctx.Value("offset")

	if err != nil {
		return nil, err
	}

	// perform query
	var todos []types.Todos

	// pagination purpose for optimization
	prepareQuery := "SELECT id, title, description, completed, due_date, created_at, updated_at FROM todos WHERE user_id = $1 LIMIT $2 OFFSET $3"

	rows, err := conn.Query(ctx, prepareQuery, uuidUserId, limit, offset)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var todo types.Todos

		err = rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed, &todo.DueDate, &todo.CreatedAt, &todo.UpdatedAt)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (s *TodosStore) Update(ctx context.Context, req types.TodosPutRequestBody) (*types.Todos, error) {
	// acquire connection
	conn, err := s.db.Acquire(ctx)

	if err != nil {
		return nil, err
	}
	// defer release connection
	defer conn.Release()

	// get user id from context
	userId := ctx.Value("user-id")
	uuidUserId, err := uuid.Parse(userId.(string))

	if err != nil {
		return nil, err
	}

	// perform query
	var todo types.Todos

	prepareQuery := "UPDATE todos SET title = $1, description = $2, completed = $3, due_date = $4 WHERE id = $5 AND user_id = $6 RETURNING id, title, description, completed, due_date, created_at, updated_at"

	err = conn.QueryRow(ctx, prepareQuery, req.Title, req.Description, req.Completed, req.DueDate, req.Id, uuidUserId).Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed, &todo.DueDate, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s *TodosStore) Delete(ctx context.Context, req types.TodosDeleteRequestBody) error {
	// acquire connection
	conn, err := s.db.Acquire(ctx)

	if err != nil {
		return err
	}
	// defer release connection
	defer conn.Release()

	// get user id from context
	userId := ctx.Value("user-id")
	uuidUserId, err := uuid.Parse(userId.(string))

	if err != nil {
		return err
	}

	// perform query
	prepareQuery := "DELETE FROM todos WHERE id = $1 AND user_id = $2"

	_, err = conn.Exec(ctx, prepareQuery, req.Id, uuidUserId)

	if err != nil {
		return err
	}

	return nil
}
