-- name: CreateUser :one
INSERT INTO users(
name,
email
) VALUES (
    $1, $2
) RETURNING *;