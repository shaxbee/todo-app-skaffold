// +build integration

package server_test

import (
	"context"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/shaxbee/todo-app-skaffold/pkg/api"
	"github.com/shaxbee/todo-app-skaffold/pkg/dbtest"
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
	addr := container.Listener().Addr().String()

	// wait for server to be listening
	waitForServer(t, addr)

	client := api.NewAPIClient(&api.Configuration{
		Servers: []api.ServerConfiguration{{
			URL: "http://" + addr,
		}},
	})

	title := "buy milk"
	content := "buy 2l of full fat milk"

	createTodo := func(t *testing.T, title string, content string) api.CreateTodoResponse {
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

	deleteAllTodos := func(t *testing.T) {
		httpRes, err := client.TodoApi.DeleteAllTodos(ctx).Execute()
		if err != nil {
			t.Fatalf("failed to delete all todos: %v", err)
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

		todo, httpRes, err := client.TodoApi.GetTodo(ctx, createRes.Id).Execute()
		if err != nil {
			t.Fatalf("failed to get todo: %v", err)
		}

		assert.Equal(t, http.StatusOK, httpRes.StatusCode, "failed to get todo")
		assert.Equal(t, api.Todo{
			Id:      createRes.Id,
			Title:   todo.Title,
			Content: todo.Content,
		}, todo)
	})

	t.Run("list todos", func(t *testing.T) {
		t.Cleanup(func() { deleteAllTodos(t) })

		createRes := createTodo(t, title, content)

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

	// shutdown container
	cancel()

	if err := errg.Wait(); err != nil {
		t.Fatalf("server error: %v", err)
	}
}

func waitForServer(t *testing.T, addr string) {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 10 * time.Millisecond
	bo.MaxElapsedTime = 5 * time.Second

	err := backoff.Retry(func() error {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		conn.Close()

		return nil
	}, bo)
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}
}
