-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.pets_to_tags (
    pet_id     bigint REFERENCES public.pets (id) ON DELETE CASCADE ON UPDATE CASCADE,
    tag_id     bigint REFERENCES public.tags (id) ON DELETE CASCADE ON UPDATE CASCADE,
    UNIQUE (pet_id, tag_id),
    created_at timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.pets_to_tags;
-- +goose StatementEnd
