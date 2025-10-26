-- name: GetById :one
SELECT * FROM public.pet_categories
WHERE id = sqlc.arg('id');
