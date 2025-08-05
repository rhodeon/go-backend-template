-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users (
  id              bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  email           text NOT NULL UNIQUE,
  username        text NOT NULL UNIQUE,
  first_name      text NOT NULL,
  last_name       text NOT NULL,
  phone_number    text,
  hashed_password text NOT NULL,
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.users;
-- +goose StatementEnd
