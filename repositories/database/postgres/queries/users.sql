-- name: Create :one
INSERT INTO public.users (
  username,
  first_name,
  last_name,
  email,
  phone_number,
  hashed_password
)
VALUES (
  sqlc.arg('username'),
  sqlc.arg('first_name'),
  sqlc.arg('last_name'),
  sqlc.arg('email'),
  sqlc.arg('phone_number'),
  sqlc.arg('hashed_password')
)
RETURNING *;

-- name: GetById :one
SELECT *
FROM public.users
WHERE id = sqlc.arg('id');

-- name: Update :one
UPDATE public.users
SET
  email = sqlc.arg('email'),
  username = sqlc.arg('username')
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: Delete :one
DELETE
FROM public.users
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: GetByEmail :one
SELECT *
FROM public.users
WHERE email = sqlc.arg('email');

-- name: GetByUsername :one
SELECT *
FROM public.users
WHERE username = sqlc.arg('username');


-- name: Verify :exec
UPDATE public.users
SET is_verified = TRUE
WHERE id = sqlc.arg('id');
