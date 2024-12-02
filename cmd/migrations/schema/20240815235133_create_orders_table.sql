-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.posts
(
    id         int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    content    text        NOT NULL,
    user_id    int REFERENCES public.users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updates_at timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.posts;
-- +goose StatementEnd
