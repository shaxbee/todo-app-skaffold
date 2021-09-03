package todo

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/shaxbee/todo-app-skaffold/api"
	"github.com/shaxbee/todo-app-skaffold/internal/httprouter"
	"github.com/shaxbee/todo-app-skaffold/services/todo/model"
)

type Server struct {
	queries *model.Queries
}

func NewServer(db model.DBTX) *Server {
	return &Server{
		queries: model.New(db),
	}
}

func (s *Server) RegisterRoutes(router *httprouter.Router) {
	router.Handler(http.MethodPost, "/api/v1/todo", s.create)
	router.Handler(http.MethodGet, "/api/v1/todo/:id", s.get)
	router.Handler(http.MethodGet, "/api/v1/todo", s.list)
	router.Handler(http.MethodDelete, "/api/v1/todo/:id", s.delete)
	router.Handler(http.MethodDelete, "/api/v1/todo", s.deleteAll)
}

func (s *Server) create(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	var ctReq api.CreateTodoRequest
	if err := httprouter.JSONRequest(req, &ctReq); err != nil {
		return err
	}

	if len(ctReq.Title) > 20 {
		return httprouter.NewError(http.StatusBadRequest, httprouter.Message("title should have maximum length of 20 characters"))
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

	return httprouter.JSONResponse(w, http.StatusCreated, api.CreateTodoResponse{
		Id: id,
	})
}

func (s *Server) get(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	rawID := httprouter.ParamsFromContext(ctx).ByName("id")

	id, err := uuid.Parse(rawID)
	if err != nil {
		return httprouter.NewError(
			http.StatusBadRequest,
			httprouter.Message("invalid id"),
			httprouter.Cause(err),
		)
	}

	t, err := s.queries.Get(ctx, id)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return httprouter.NewError(
			http.StatusNotFound,
			httprouter.Messagef("todo %q not found", id),
			httprouter.Operational(),
		)
	case err != nil:
		return fmt.Errorf("failed to get todo: %w", err)
	}

	return httprouter.JSONResponse(w, http.StatusOK, api.Todo{
		Id:      t.ID,
		Title:   t.Title,
		Content: t.Content,
	})
}

func (s *Server) list(w http.ResponseWriter, req *http.Request) error {
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

	return httprouter.JSONResponse(w, http.StatusOK, resTodos)
}

func (s *Server) delete(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	rawID := httprouter.ParamsFromContext(ctx).ByName("id")

	id, err := uuid.Parse(rawID)
	if err != nil {
		return httprouter.NewError(
			http.StatusBadRequest,
			httprouter.Message("invalid id"),
			httprouter.Cause(err),
		)
	}

	n, err := s.queries.Delete(ctx, id)
	switch {
	case err != nil:
		return fmt.Errorf("failed to delete todo: %w", err)
	case n == 0:
		return httprouter.NewError(
			http.StatusNotFound,
			httprouter.Messagef("todo %q not found", id),
			httprouter.Operational(),
		)
	default:
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}

func (s *Server) deleteAll(w http.ResponseWriter, req *http.Request) error {
	if err := s.queries.DeleteAll(req.Context()); err != nil {
		return fmt.Errorf("failed to delete all todos: %w", err)
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
