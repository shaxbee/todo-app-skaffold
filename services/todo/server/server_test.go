// +build integration

package server_test

import (
	"context"
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

	createTodo := func(t *testing.T, title string, content string) uuid.UUID {
		//nolint:bodyclose
		res, httpRes, err := client.TodoApi.CreateTodo(ctx).CreateTodoRequest(api.CreateTodoRequest{
			Title:   title,
			Content: content,
		}).Execute()
		if err != nil {
			t.Fatalf("failed to create todo: %v", err)
		}

		assert.Equal(t, http.StatusCreated, httpRes.StatusCode, "failed to create todo: unexpected status")

		return res.Id
	}

	getTodo := func(t *testing.T, id uuid.UUID) (api.Todo, bool) {
		//nolint:bodyclose
		res, httpRes, err := client.TodoApi.GetTodo(ctx, id).Execute()
		switch {
		case err != nil && httpRes == nil:
			t.Errorf("failed to get todo %q: %v", id, err)
			return res, false
		case err != nil && httpRes.StatusCode != http.StatusOK && httpRes.StatusCode != http.StatusNotFound:
			t.Errorf("failed to get todo %q: unexpected status code %d", id, httpRes.StatusCode)
			return res, false
		default:
			return res, httpRes.StatusCode == http.StatusOK
		}
	}

	deleteTodo := func(t *testing.T, id uuid.UUID) bool {
		//nolint:bodyclose
		httpRes, err := client.TodoApi.DeleteTodo(ctx, id).Execute()
		switch {
		case err != nil && httpRes == nil:
			t.Errorf("failed to delete todo %q: %v", id, err)
			return false
		case err != nil && httpRes.StatusCode != http.StatusNoContent && httpRes.StatusCode != http.StatusNotFound:
			t.Errorf("failed to delete todo %q: unexpected status code %d", id, httpRes.StatusCode)
			return false
		default:
			return httpRes.StatusCode == http.StatusNoContent
		}
	}

	t.Run("create todo", func(t *testing.T) {
		id := createTodo(t, title, content)
		t.Cleanup(func() { deleteTodo(t, id) })

		assert.NotZero(t, id, "expected non zero id")

		otherID := createTodo(t, title, content)
		t.Cleanup(func() { deleteTodo(t, otherID) })

		assert.NotEqual(t, otherID, id, "expected unique id")
	})

	t.Run("get todo", func(t *testing.T) {
		id := createTodo(t, title, content)
		t.Cleanup(func() { deleteTodo(t, id) })

		todo, exists := getTodo(t, id)

		assert.True(t, exists, "expected todo to exist")
		assert.Equal(t, api.Todo{
			Id:      id,
			Title:   title,
			Content: content,
		}, todo)

		_, exists = getTodo(t, uuid.New())
		assert.False(t, exists)
	})

	t.Run("list todos", func(t *testing.T) {
		id := createTodo(t, title, content)

		//nolint:bodyclose
		todos, httpRes, err := client.TodoApi.ListTodos(ctx).Execute()
		if err != nil {
			t.Fatalf("failed to list todos: %v", err)
		}

		assert.Equal(t, http.StatusOK, httpRes.StatusCode, "failed to list todos: unexpected status")
		assert.Contains(t, todos, api.Todo{
			Id:      id,
			Title:   title,
			Content: content,
		}, "expected todo to match")
	})

	t.Run("delete todo", func(t *testing.T) {
		id := createTodo(t, title, content)

		assert.True(t, deleteTodo(t, id), "expected todo to be deleted")

		_, exists := getTodo(t, id)
		assert.False(t, exists, "expected todo to be deleted")

		assert.False(t, deleteTodo(t, id), "expected previously deleted todo to be not found")
		assert.False(t, deleteTodo(t, uuid.New()), "expected non-existent todo to be not found")
	})

	// shutdown container
	cancel()

	if err := errg.Wait(); err != nil {
		t.Fatalf("server error: %v", err)
	}
}
