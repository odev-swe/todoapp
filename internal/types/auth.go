package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserIdKey string

type User struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Token     Token     `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Token struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type UserRequestBody struct {
	Email    string `json:"email" example:"admin@gmail.com"`
	Password string `json:"password" example:"admin"`
}

type AuthServices interface {
	Register(ctx context.Context, user UserRequestBody) (*User, error)
	Login(ctx context.Context, user UserRequestBody) (*User, error)
}
