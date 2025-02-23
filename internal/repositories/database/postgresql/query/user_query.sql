-- name: CreateUser :one
INSERT INTO users (username,
                   password)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: GetUserForLogin :one
SELECT *
FROM users
WHERE username = $1
  and password = $2;

-- name: GetUsersById :one
SELECT *
From users
WHERE id = $1;

