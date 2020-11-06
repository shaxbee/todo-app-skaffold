// +build integration

package server_test

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/shaxbee/todo-app-skaffold/pkg/api"
	"github.com/shaxbee/todo-app-skaffold/pkg/dbtest"
	"github.com/shaxbee/todo-app-skaffold/pkg/servertest"
	"github.com/shaxbee/todo-app-skaffold/services/todo/server"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestAPI(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// parse config overriding address to random localhost port
	config, err := server.ParseConfig(server.ConfigAddr("127.0.0.1:0"))
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	enabled := os.Getenv("PGHOST") == "" && os.Getenv("PGPORT") == ""

	db := dbtest.Postgres(
		t,
		dbtest.Migration("../migrations"),
		dbtest.Enabled(enabled),
	)

	// setup server container
	container := server.NewContainer(config, server.ContainerDB(db))

	errg := container.Run(ctx)
	addr := container.Addr()

	// wait for server to be listening
	servertest.WaitForServer(t, servertest.Addr(addr))

	client := api.NewAPIClient(&api.Configuration{
		Servers: []api.ServerConfiguration{{
			URL: "http://" + addr,
		}},
	})

	title := "buy milk"
	content := "buy 2l of full fat milk"

	createTodo := func(t *testing.T, title string, content string) api.CreateTodoResponse {
		//nolint:bodyclose
		res, httpRes, err := client.TodoApi.CreateTodo(ctx).CreateTodoRequest(api.CreateTodoRequest{
			Title:   title,
			Content: content,
		}).Execute()
		if err != nil {
			t.Fatalf("failed to create todo: %v", err)
		}

		assert.Equal(t, http.StatusCreated, httpRes.StatusCode, "failed to create todo")

		return res
	}

	getTodo := func(t *testing.T, id uuid.UUID) api.Todo {
		//nolint:bodyclose
		res, httpRes, err := client.TodoApi.GetTodo(ctx, id).Execute()
		if err != nil {
			t.Fatalf("failed to get todo %q: %v", id, err)
		}

		assert.Equal(t, http.StatusOK, httpRes.StatusCode, "failed to get todo")

		return res
	}

	todoExists := func(t *testing.T, id uuid.UUID) bool {
		//nolint:bodyclose
		_, httpRes, err := client.TodoApi.GetTodo(ctx, id).Execute()

		var apiErr api.GenericOpenAPIError
		if err != nil && !errors.As(err, &apiErr) {
			t.Fatalf("failed to get todo %q: %v", id, err)
		}

		switch httpRes.StatusCode {
		case http.StatusOK:
			return true
		case http.StatusNotFound:
			return false
		default:
			t.Fatalf("failed to get todo %q: %v", id, err)
			return false
		}
	}

	deleteAllTodos := func(t *testing.T) {
		//nolint:bodyclose
		httpRes, err := client.TodoApi.DeleteAllTodos(ctx).Execute()
		if err != nil {
			t.Errorf("failed to delete all todos: %v", err)
		}

		assert.Equal(t, http.StatusNoContent, httpRes.StatusCode, "failed to delete all todos")
	}

	t.Run("create todo", func(t *testing.T) {
		t.Cleanup(func() { deleteAllTodos(t) })

		res := createTodo(t, title, content)
		assert.NotZero(t, res.Id, "expected non zero id")
	})

	t.Run("get todo", func(t *testing.T) {
		t.Cleanup(func() { deleteAllTodos(t) })

		createRes := createTodo(t, title, content)
		todo := getTodo(t, createRes.Id)

		assert.Equal(t, api.Todo{
			Id:      createRes.Id,
			Title:   todo.Title,
			Content: todo.Content,
		}, todo)
		assert.False(t, todoExists(t, uuid.New()))
	})

	t.Run("list todos", func(t *testing.T) {
		t.Cleanup(func() { deleteAllTodos(t) })

		createRes := createTodo(t, title, content)

		//nolint:bodyclose
		todos, httpRes, err := client.TodoApi.ListTodos(ctx).Execute()
		if err != nil {
			t.Fatalf("failed to list todos: %v", err)
		}

		assert.Equal(t, http.StatusOK, httpRes.StatusCode, "failed to list todos")
		assert.Equal(t, []api.Todo{
			{
				Id:      createRes.Id,
				Title:   title,
				Content: content,
			},
		}, todos)
	})

	t.Run("delete todo", func(t *testing.T) {
		t.Cleanup(func() { deleteAllTodos(t) })

		createRes := createTodo(t, title, content)

		//nolint:bodyclose
		httpRes, err := client.TodoApi.DeleteTodo(ctx, createRes.Id).Execute()
		if err != nil {
			t.Fatalf("failed to delete todo: %v", err)
		}

		assert.Equal(t, http.StatusNoContent, httpRes.StatusCode, "failed to delete todo")
		assert.False(t, todoExists(t, createRes.Id))
	})

	// shutdown container
	cancel()

	if err := errg.Wait(); err != nil {
		t.Fatalf("server error: %v", err)
	}
}
