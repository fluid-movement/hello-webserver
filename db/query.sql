-- name: ListTodos :many
SELECT id, title, done, created_at FROM todos ORDER BY created_at ASC;

-- name: CreateTodo :one
INSERT INTO todos (title) VALUES ($1)
RETURNING id, title, done, created_at;

-- name: ToggleTodo :one
UPDATE todos SET done = NOT done WHERE id = $1
RETURNING id, title, done, created_at;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1;
