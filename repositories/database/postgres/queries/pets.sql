-- name: Create :one
INSERT INTO public.pets (name, category_id, image_urls)
VALUES (@name, @category_id, @image_urls)
RETURNING *;
