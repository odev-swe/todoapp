package services

import (
	"context"

	"github.com/odev-swe/todoapp/internal/store"
	"github.com/odev-swe/todoapp/internal/types"
)

type AuthService struct {
	store *store.AuthStore
}

func NewAuthService(store *store.AuthStore) *AuthService {
	return &AuthService{store: store}
}

func (s *AuthService) Register(ctx context.Context, user types.UserRequestBody) (*types.User, error) {
	return s.store.Register(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, user types.UserRequestBody) (*types.User, error) {
	return s.store.Login(ctx, user)
}
