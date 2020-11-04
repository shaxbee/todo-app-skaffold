package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/shaxbee/todo-app-skaffold/pkg/api/todo"
	"github.com/shaxbee/todo-app-skaffold/pkg/httperror"
	"github.com/shaxbee/todo-app-skaffold/services/todo/model"
)

type TodoServer struct {
	queries *model.Queries
}

func New(db *sql.DB) *TodoServer {
	return &TodoServer{
		queries: model.New(db),
	}
}

func (s *TodoServer) RegisterRoutes(router *httprouter.Router, errorMiddleware httperror.Middleware) {
	router.Handler(http.MethodGet, "/api/v1/todo/:id", errorMiddleware(s.Get))
	router.Handler(http.MethodGet, "/api/v1/todo", errorMiddleware(s.List))
	router.Handler(http.MethodPost, "/api/v1/todo", errorMiddleware(s.Create))
}

func (s *TodoServer) Get(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	rawId := httprouter.ParamsFromContext(ctx).ByName("id")
	id, err := uuid.Parse(rawId)
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

	data, err := json.Marshal(todo.Todo{
		Id:      t.ID,
		Title:   t.Title,
		Content: t.Content,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal Todo: %w", err)
	}

	_, _ = w.Write(data)
	w.WriteHeader(http.StatusOK)

	return nil
}

func (s *TodoServer) List(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	todos, err := s.queries.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list todos: %w", err)
	}

	resTodos := make([]todo.Todo, len(todos))
	for i, t := range todos {
		resTodos[i] = todo.Todo{
			Id:      t.ID,
			Title:   t.Title,
			Content: t.Content,
		}
	}

	data, err := json.Marshal(resTodos)
	if err != nil {
		return err
	}

	_, _ = w.Write(data)
	return nil
}

func (s *TodoServer) Create(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()

	reqData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	var ctReq todo.CreateTodoRequest
	if err := json.Unmarshal(reqData, &ctReq); err != nil {
		return httperror.New(
			http.StatusBadRequest,
			httperror.Cause(err),
		)
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

	resData, err := json.Marshal(todo.CreateTodoResponse{
		Id: id,
	})
	if err != nil {
		return err
	}

	_, _ = w.Write(resData)
	return nil
}
