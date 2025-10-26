-- name: Create :one
INSERT INTO public.pets (name, category_id, image_urls)
VALUES (sqlc.arg('name'), sqlc.arg('category_id'), sqlc.arg('image_urls'))
RETURNING *;
