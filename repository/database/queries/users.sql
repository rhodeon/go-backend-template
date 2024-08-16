-- name: Create :one
INSERT INTO users(email, username)
VALUES (@email, @username)
RETURNING *;

-- name: GetById :one
SELECT *
FROM users
WHERE id = @id;

-- name: Update :one
UPDATE users
SET email    = @email,
    username = @username
WHERE id = @id
RETURNING *;

-- name: Delete :one
DELETE
FROM users
WHERE id = @id
RETURNING *;

-- name: ListPosts :many
SELECT posts.*
FROM users
         JOIN posts on users.id = posts.user_id;