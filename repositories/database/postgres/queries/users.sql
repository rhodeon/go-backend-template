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
  @username,
  @first_name,
  @last_name,
  @email,
  @phone_number,
  @hashed_password
)
RETURNING *;

-- name: GetById :one
SELECT *
FROM public.users
WHERE id = @id;

-- name: Update :one
UPDATE public.users
SET
  email = @email,
  username = @username
WHERE id = @id
RETURNING *;

-- name: Delete :one
DELETE
FROM public.users
WHERE id = @id
RETURNING *;

-- name: GetByEmail :one
SELECT *
FROM public.users
WHERE email = @email;

-- name: GetByUsername :one
SELECT *
FROM public.users
WHERE username = @username;


-- name: Verify :exec
UPDATE public.users
SET is_verified = TRUE
WHERE id = @id;
