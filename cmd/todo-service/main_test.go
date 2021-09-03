//go:build integration
// +build integration

package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/shaxbee/todo-app-skaffold/api"
	"github.com/shaxbee/todo-app-skaffold/internal/dbtest"
	"github.com/shaxbee/todo-app-skaffold/internal/servertest"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestAPI(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := servertest.Setup(t, servertest.MakeHandler(func() http.Handler {
		config, err := parseConfig()
		if err != nil {
			t.Fatal(err)
		}

		config.Dev = true

		cont := newContainer(config)
		cont.state.db = dbtest.SetupPostgres(t, dbtest.Migration("../../services/todo/migrations"))

		return cont.httpRouter()
	}))

	client := api.NewAPIClient(&api.Configuration{
		Servers: []api.ServerConfiguration{{
			URL: endpoint,
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

		if http.StatusCreated != httpRes.StatusCode {
			t.Errorf("failed to create todo: unexpected status %d", httpRes.StatusCode)
		}

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

	deleteAllTodos := func(t *testing.T) bool {
		//nolint:bodyclose
		httpRes, err := client.TodoApi.DeleteAllTodos(ctx).Execute()
		switch {
		case err != nil && httpRes == nil:
			t.Errorf("failed to delete all todos: %v", err)
			return false
		case err != nil && httpRes.StatusCode != http.StatusNoContent:
			t.Errorf("failed to delete all todos: unexpected status code %d", httpRes.StatusCode)
			return false
		default:
			return httpRes.StatusCode == http.StatusNoContent
		}
	}

	t.Run("create todo", func(t *testing.T) {
		id := createTodo(t, title, content)
		t.Cleanup(func() { deleteTodo(t, id) })

		if id == uuid.Nil {
			t.Error("expected non zero id")
		}

		otherID := createTodo(t, title, content)
		t.Cleanup(func() { deleteTodo(t, otherID) })

		if otherID == id {
			t.Error("expected unique id")
		}
	})

	t.Run("get todo", func(t *testing.T) {
		id := createTodo(t, title, content)
		t.Cleanup(func() { deleteTodo(t, id) })

		actual, exists := getTodo(t, id)

		if !exists {
			t.Error("expected todo to exist")
		}

		expected := api.Todo{
			Id:      id,
			Title:   title,
			Content: content,
		}

		if diff := cmp.Diff(expected, actual); diff != "" {
			t.Error("expected equal todo:", diff)
		}

		if _, exists = getTodo(t, uuid.New()); exists {
			t.Error("unexpected todo")
		}
	})

	t.Run("list todos", func(t *testing.T) {
		if !deleteAllTodos(t) {
			t.FailNow()
		}

		id := createTodo(t, title, content)

		//nolint:bodyclose
		actual, httpRes, err := client.TodoApi.ListTodos(ctx).Execute()
		if err != nil {
			t.Fatalf("failed to list todos: %v", err)
		}

		if http.StatusOK != httpRes.StatusCode {
			t.Errorf("failed to list todos: unexpected status %d", httpRes.StatusCode)
		}

		expected := []api.Todo{{
			Id:      id,
			Title:   title,
			Content: content,
		}}

		if diff := cmp.Diff(expected, actual); diff != "" {
			t.Error("expected equal todos:", diff)
		}
	})

	t.Run("delete todo", func(t *testing.T) {
		id := createTodo(t, title, content)

		if !deleteTodo(t, id) {
			t.Error("expected todo to be deleted")
		}

		if _, exists := getTodo(t, id); exists {
			t.Error("expected todo to be deleted")
		}

		if deleteTodo(t, id) {
			t.Error("expected previously deleted todo to be not found")
		}

		if deleteTodo(t, uuid.New()) {
			t.Error("expected non-existent todo to be not found")
		}
	})
}
