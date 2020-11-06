package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/shaxbee/todo-app-skaffold/pkg/api"
	"github.com/shaxbee/todo-app-skaffold/pkg/httperror"
	"github.com/shaxbee/todo-app-skaffold/pkg/routes"
	"github.com/shaxbee/todo-app-skaffold/services/todo/model"
)

type TodoServer struct {
	queries *model.Queries
}

func New(db model.DBTX) *TodoServer {
	return &TodoServer{
		queries: model.New(db),
	}
}

func (s *TodoServer) RegisterRoutes(router *httprouter.Router, errorMiddleware httperror.Middleware) {
	router.Handler(http.MethodGet, "/api/v1/todo/:id", errorMiddleware(s.Get))
	router.Handler(http.MethodGet, "/api/v1/todo", errorMiddleware(s.List))
	router.Handler(http.MethodPost, "/api/v1/todo", errorMiddleware(s.Create))
	router.Handler(http.MethodDelete, "/api/v1/todo", errorMiddleware(s.DeleteAll))
}

func (s *TodoServer) Get(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	rawID := httprouter.ParamsFromContext(ctx).ByName("id")

	id, err := uuid.Parse(rawID)
	if err != nil {
		return httperror.New(http.StatusBadRequest, httperror.Message("invalid id"), httperror.Cause(err))
	}

	t, err := s.queries.Get(ctx, id)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return httperror.New(http.StatusNotFound, httperror.Message("todo not found"))
	case err != nil:
		return fmt.Errorf("failed to get todo: %w", err)
	}

	return routes.JSONResponseBody(w, http.StatusOK, api.Todo{
		Id:      t.ID,
		Title:   t.Title,
		Content: t.Content,
	})
}

func (s *TodoServer) List(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	todos, err := s.queries.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list todos: %w", err)
	}

	resTodos := make([]api.Todo, len(todos))
	for i, t := range todos {
		resTodos[i] = api.Todo{
			Id:      t.ID,
			Title:   t.Title,
			Content: t.Content,
		}
	}

	return routes.JSONResponseBody(w, http.StatusOK, resTodos)
}

func (s *TodoServer) Create(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	var ctReq api.CreateTodoRequest
	if err := routes.JSONRequestBody(req, &ctReq); err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate todo id: %w", err)
	}

	err = s.queries.Create(ctx, model.CreateParams{
		ID:      id,
		Title:   ctReq.Title,
		Content: ctReq.Content,
	})
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	return routes.JSONResponseBody(w, http.StatusCreated, api.CreateTodoResponse{
		Id: id,
	})
}

func (s *TodoServer) DeleteAll(w http.ResponseWriter, req *http.Request) error {
	if err := s.queries.DeleteAll(req.Context()); err != nil {
		return fmt.Errorf("failed to delete all todos: %w", err)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
