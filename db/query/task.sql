-- name: CreateTask :one
INSERT INTO tasks (owner,title,description)
VALUES ($1,$2,$3)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1 
LIMIT 1;

-- name: ListTasks :many
SELECT * FROM tasks
WHERE owner = $1
ORDER BY id
OFFSET $2
LIMIT $3;

-- name: UpdateTask :one
UPDATE tasks
SET title = $2,
description = $3,
done = $4
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;