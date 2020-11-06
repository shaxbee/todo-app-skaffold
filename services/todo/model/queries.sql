-- name: Get :one
SELECT * FROM todo WHERE id=sqlc.arg(id);

-- name: List :many
SELECT * FROM todo;

-- name: Create :exec
INSERT INTO todo (id, title, content) VALUES (sqlc.arg(id), sqlc.arg(title), sqlc.arg(content));

-- name: Delete :exec
DELETE FROM todo WHERE id=sqlc.arg(id);

-- name: DeleteAll :exec
DELETE FROM todo;