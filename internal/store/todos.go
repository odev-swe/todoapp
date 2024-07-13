package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/odev-swe/todoapp/internal/types"
	"github.com/odev-swe/todoapp/libs"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type TodosStore struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewTodosStore(db *pgxpool.Pool, redis *redis.Client) *TodosStore {
	return &TodosStore{
		db:    db,
		redis: redis,
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
	userId := ctx.Value(types.UserIdKey("user-id"))
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

	err = deleteCache(ctx, "todos", s.redis)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s *TodosStore) Get(ctx context.Context) ([]types.Todos, error) {
	var todo types.Todos
	var todos []types.Todos

	// check cache first
	jsonData, err := s.redis.Get(ctx, "todos").Result()

	if err == redis.Nil {
		// acquire connection
		conn, err := s.db.Acquire(ctx)

		if err != nil {
			return nil, err
		}

		// defer release connection
		defer conn.Release()

		// get user id from context
		userId := ctx.Value(types.UserIdKey("user-id"))
		uuidUserId, err := uuid.Parse(userId.(string))

		limit := ctx.Value("limit")
		offset := ctx.Value("offset")

		if err != nil {
			return nil, err
		}

		// perform query

		// pagination purpose for optimization
		prepareQuery := "SELECT id, title, description, completed, due_date, created_at, updated_at FROM todos WHERE user_id = $1 LIMIT $2 OFFSET $3"

		rows, err := conn.Query(ctx, prepareQuery, uuidUserId, limit, offset)

		if err != nil {
			return nil, err
		}

		for rows.Next() {

			err = rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed, &todo.DueDate, &todo.CreatedAt, &todo.UpdatedAt)

			if err != nil {
				return nil, err
			}

			todos = append(todos, todo)
		}

		// set cache
		data, err := libs.StringifyJSON(todos)

		if err != nil {
			return nil, err
		}

		// expiration in 10 seconds
		zap.L().Info("retrieve data from database")
		err = s.redis.Set(ctx, "todos", data, 30*time.Second).Err()

		if err != nil {
			return nil, err
		}

		return todos, nil
	}

	zap.L().Info("retrieve data from cache")

	err = libs.ParseStringJSON(jsonData, &todos)

	if err != nil {
		return nil, err
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
	userId := ctx.Value(types.UserIdKey("user-id"))
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

	err = deleteCache(ctx, "todos", s.redis)

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
	userId := ctx.Value(types.UserIdKey("user-id"))
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

func deleteCache(ctx context.Context, key string, r *redis.Client) error {
	_, err := r.Get(ctx, key).Result()

	if err == redis.Nil {
		return nil
	}

	err = r.Del(ctx, key).Err()

	if err != nil {
		return err
	}

	return nil
}
