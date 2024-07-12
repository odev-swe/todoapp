package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/odev-swe/todoapp/internal/types"
	"github.com/odev-swe/todoapp/libs"
	"golang.org/x/crypto/bcrypt"
)

type AuthStore struct {
	db     *pgxpool.Pool
	secret string
}

func NewAuthStore(db *pgxpool.Pool, secret string) *AuthStore {
	return &AuthStore{
		db:     db,
		secret: secret,
	}
}

func (s *AuthStore) Register(ctx context.Context, req types.UserRequestBody) (*types.User, error) {
	// Acquire a connection from the pool
	conn, err := s.db.Acquire(ctx)

	// Release the connection back to the pool
	defer conn.Release()

	if err != nil {
		// handle error
		return nil, err
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	var user types.User

	// perform the query
	// register query statement
	prepareQuery := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, created_at, updated_at"

	err = conn.QueryRow(ctx, prepareQuery, req.Email, hashedPassword).Scan(&user.Id, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthStore) Login(ctx context.Context, req types.UserRequestBody) (*types.User, error) {
	// Acquire a connection from the pool
	conn, err := s.db.Acquire(ctx)

	// Release the connection back to the pool
	defer conn.Release()

	if err != nil {
		// handle error
		return nil, err
	}

	var user types.User

	// perform the query
	// login query statement
	prepareQuery := "SELECT id, email,password, created_at, updated_at FROM users WHERE email = $1"

	err = conn.QueryRow(ctx, prepareQuery, req.Email).Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		return nil, err
	}

	// clear the password
	user.Password = ""

	// generate access token
	at, err := libs.GenerateToken(user, s.secret, libs.AccessToken)

	if err != nil {
		return nil, err
	}

	// generate refresh token
	rt, err := libs.GenerateToken(user, s.secret, libs.RefreshToken)

	if err != nil {
		return nil, err
	}

	user.Token.AccessToken = at
	user.Token.RefreshToken = rt

	return &user, nil
}
