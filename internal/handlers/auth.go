package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/odev-swe/todoapp/internal/types"
	"github.com/odev-swe/todoapp/libs"
)

type AuthHandler struct {
	authService types.AuthServices
}

func NewAuthHandler(authService types.AuthServices) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterRoute(r chi.Router) {
	// handle the request
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
}

// Auth godoc
//	@Summary		Register an account
//	@Description	register an account with email and password
//	@Tags			auth	
//	@Accept			json
//	@Produce		json
//	@Param			body	body		types.UserRequestBody	true	"User object that needs to be registered"
//	@Success		200		{object}	libs.Response
//	@Failure		400		{object}	libs.Response
//	@Failure		500		{object}	libs.Response
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// handle the request
	var user types.UserRequestBody

	err := libs.ParseJSON(r, &user)

	if err != nil {
		libs.BadRequest(w, "Invalid request body")
		return
	}

	// timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.authService.Register(ctx, user)

	if err != nil {
		// handle context timeout error
		if ctx.Err() == context.DeadlineExceeded {
			libs.InternalServerError(w, "Request timeout")
			return
		}
		libs.InternalServerError(w, err.Error())
		return
	}

	libs.WriteJSON(w, true, http.StatusCreated, "User created successfully", res)
}

// Auth godoc
//	@Summary		Login an account
//	@Description	Login an account with email and password
//	@Tags			auth	
//	@Accept			json
//	@Produce		json
//	@Param			body	body		types.UserRequestBody	true	"User object that needs to be registered"
//	@Success		200		{object}	libs.Response
//	@Failure		400		{object}	libs.Response
//	@Failure		500		{object}	libs.Response
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// handle the request
	var user types.UserRequestBody

	err := libs.ParseJSON(r, &user)

	if err != nil {
		libs.BadRequest(w, "Invalid request body")
		return
	}

	// timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.authService.Login(ctx, user)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			libs.InternalServerError(w, "Request timeout")
			return
		}
		libs.InternalServerError(w, err.Error())
		return
	}

	libs.WriteJSON(w, true, http.StatusOK, "User logged in successfully", res)
}
