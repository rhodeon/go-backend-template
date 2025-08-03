-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.pet_categories (
    id   bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.pet_categories;
-- +goose StatementEnd
