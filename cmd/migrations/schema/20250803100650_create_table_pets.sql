-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.pets (
  id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name        text NOT NULL,
  category_id bigint REFERENCES public.pet_categories (id),
  status      text CHECK (status IN ('available', 'pending', 'sold')),
  image_urls  text [] NOT NULL DEFAULT '{}',
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.pets;
-- +goose StatementEnd
