package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/odev-swe/todoapp/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of the AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, user types.UserRequestBody) (*types.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, user types.UserRequestBody) (*types.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*types.User), args.Error(1)
}

const UUIDtest = "3162d3f0-5532-402d-ab85-28946a279cac"

func TestRegister(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	tests := []struct {
		name           string
		inputJSON      string
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Successful Registration",
			inputJSON: `{"email":"test@example.com","password":"password123"}`,
			mockBehavior: func() {
				mockService.On("Register", mock.Anything, types.UserRequestBody{
					Email:    "test@example.com",
					Password: "password123",
				}).Return(&types.User{Id: uuid.MustParse(UUIDtest), Email: "test@example.com"}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   fmt.Sprintf(`{"success":true,"message":"User created successfully","data":{"id":"%s","email":"test@example.com", "token":"{}","created_at":"0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}}`, UUIDtest),
		},
		{
			name:           "Invalid JSON",
			inputJSON:      `{"email":"test@example.com","password":}`,
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"success":false,"message":"Invalid request body"}`,
		},
		{
			name:      "Service Error",
			inputJSON: `{"email":"test@example.com","password":"password123"}`,
			mockBehavior: func() {
				mockService.On("Register", mock.Anything, types.UserRequestBody{
					Email:    "test@example.com",
					Password: "password123",
				}).Return((*types.User)(nil), errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"success":false,"message":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(tt.inputJSON))
			rr := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/register", handler.Register)
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestLogin(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	tests := []struct {
		name           string
		inputJSON      string
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Successful Login",
			inputJSON: `{"email":"test@example.com","password":"password123"}`,
			mockBehavior: func() {
				mockService.On("Login", mock.Anything, types.UserRequestBody{
					Email:    "test@example.com",
					Password: "password123",
				}).Return(&types.User{Id: uuid.New(), Email: "test@example.com", Token: types.Token{AccessToken: "some-token"}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"success":true,"message":"User logged in successfully","data":{"id":1,"email":"test@example.com","token":{"access_token":"some-token"}}}`,
		},
		{
			name:           "Invalid JSON",
			inputJSON:      `{"email":"test@example.com","password":}`,
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"success":false,"message":"Invalid request body"}`,
		},
		{
			name:      "Service Error",
			inputJSON: `{"email":"test@example.com","password":"password123"}`,
			mockBehavior: func() {
				mockService.On("Login", mock.Anything, types.UserRequestBody{
					Email:    "test@example.com",
					Password: "password123",
				}).Return((*types.User)(nil), errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"success":false,"message":"invalid credentials"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(tt.inputJSON))
			rr := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/login", handler.Login)
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestContextTimeout(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	tests := []struct {
		name           string
		endpoint       string
		inputJSON      string
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Register Timeout",
			endpoint:  "/register",
			inputJSON: `{"email":"test@example.com","password":"password123"}`,
			mockBehavior: func() {
				mockService.On("Register", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						time.Sleep(6 * time.Second)
					}).
					Return((*types.User)(nil), context.DeadlineExceeded)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"success":false,"message":"Request timeout"}`,
		},
		{
			name:      "Login Timeout",
			endpoint:  "/login",
			inputJSON: `{"email":"test@example.com","password":"password123"}`,
			mockBehavior: func() {
				mockService.On("Login", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						time.Sleep(6 * time.Second)
					}).
					Return((*types.User)(nil), context.DeadlineExceeded)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"success":false,"message":"Request timeout"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req, _ := http.NewRequest("POST", tt.endpoint, bytes.NewBufferString(tt.inputJSON))
			rr := httptest.NewRecorder()

			r := chi.NewRouter()
			if tt.endpoint == "/register" {
				r.Post("/register", handler.Register)
			} else {
				r.Post("/login", handler.Login)
			}
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}
