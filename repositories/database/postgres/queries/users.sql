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
FROM users
WHERE id = @id;

-- name: Update :one
UPDATE users
SET
  email = @email,
  username = @username
WHERE id = @id
RETURNING *;

-- name: Delete :one
DELETE
FROM users
WHERE id = @id
RETURNING *;
