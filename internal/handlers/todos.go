package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/odev-swe/todoapp/internal/types"
	"github.com/odev-swe/todoapp/libs"
)

type TodosHandler struct {
	service types.TodosServices
}

type LimitKey string
type OffsetKey string

func NewTodosHandler(service types.TodosServices) *TodosHandler {
	return &TodosHandler{service: service}
}

func (h *TodosHandler) RegisterRoute(r chi.Router) {
	// handle the request
	// post request with limit and offset parameter
	r.Get("/", h.Get)
	r.Post("/", h.Create)
	r.Put("/", h.Update)
	r.Delete("/", h.Delete)
}

// Todos godoc
//
//	@Summary		Get todos
//	@Description	get todos with pagination
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			limit	query	int	false	"Limit"		default(10)
//	@Param			offset	query	int	false	"Offset"	default(0)
//	@Security		ApiKeyAuth
//	@Success		200	{object}	libs.Response
//	@Failure		400	{object}	libs.Response
//	@Failure		500	{object}	libs.Response
//	@Router			/todos [get]
func (h *TodosHandler) Get(w http.ResponseWriter, r *http.Request) {

	// set default pagination
	var limit, offset int = 10, 10

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, LimitKey("limit"), limit)
	ctx = context.WithValue(ctx, OffsetKey("offset"), offset)

	// timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := h.service.Get(ctx)

	if err != nil {
		// handle context timeout error
		if ctx.Err() == context.DeadlineExceeded {
			libs.InternalServerError(w, "Request timeout")
			return
		}
		libs.InternalServerError(w, err.Error())
		return
	}

	libs.WriteJSON(w, true, http.StatusOK, "Todos retrieved successfully", res)
}

// Todos godoc
//
//	@Summary		Create a new todo
//	@Description	create a new todo
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//
// @Param			body	body	types.TodosPostRequestBody	true	"Todo object that needs to be created"
//
//	@Security		ApiKeyAuth
//	@Success		200	{object}	libs.Response
//	@Failure		400	{object}	libs.Response
//	@Failure		500	{object}	libs.Response
//	@Router			/todos [post]
func (h *TodosHandler) Create(w http.ResponseWriter, r *http.Request) {
	// handle the request
	var todo types.TodosPostRequestBody

	err := libs.ParseJSON(r, &todo)

	if err != nil {
		libs.BadRequest(w, "Invalid request body")
		return
	}

	// timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Create(ctx, todo)

	if err != nil {
		// handle context timeout error
		if ctx.Err() == context.DeadlineExceeded {
			libs.InternalServerError(w, "Request timeout")
			return
		}
		libs.InternalServerError(w, err.Error())
		return
	}

	libs.WriteJSON(w, true, http.StatusCreated, "Todo created successfully", res)
}

// Todos godoc
//
//	@Summary		Update a todo
//	@Description	update a todo
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			body	body	types.TodosPutRequestBody	true	"Todo object that needs to be updated"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	libs.Response
//	@Failure		400	{object}	libs.Response
//	@Failure		500	{object}	libs.Response
//	@Router			/todos [put]
func (h *TodosHandler) Update(w http.ResponseWriter, r *http.Request) {
	// handle the request
	var todo types.TodosPutRequestBody

	err := libs.ParseJSON(r, &todo)

	if err != nil {
		libs.BadRequest(w, "Invalid request body")
		return
	}

	// timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.service.Update(ctx, todo)

	if err != nil {
		// handle context timeout error
		if ctx.Err() == context.DeadlineExceeded {
			libs.InternalServerError(w, "Request timeout")
			return
		}
		libs.InternalServerError(w, err.Error())
		return
	}

	libs.WriteJSON(w, true, http.StatusOK, "Todo updated successfully", res)
}

// Todos godoc
//
//	@Summary		Delete a todo
//	@Description	delete a todo
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			body	body	types.TodosDeleteRequestBody	true	"Todo object that needs to be created"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	libs.Response
//	@Failure		400	{object}	libs.Response
//	@Failure		500	{object}	libs.Response
//	@Router			/todos [delete]
func (h *TodosHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// handle the request
	var todo types.TodosDeleteRequestBody

	err := libs.ParseJSON(r, &todo)

	if err != nil {
		libs.BadRequest(w, "Invalid request body")
		return
	}

	// timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.service.Delete(ctx, todo)

	if err != nil {
		// handle context timeout error
		if ctx.Err() == context.DeadlineExceeded {
			libs.InternalServerError(w, "Request timeout")
			return
		}
		libs.InternalServerError(w, err.Error())
		return
	}

	libs.WriteJSON(w, true, http.StatusOK, "Todo deleted successfully", nil)
}
