-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.orders (
  id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

  user_id       bigint NOT NULL REFERENCES public.users (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE,

  pet_id        bigint NOT NULL REFERENCES public.pets (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE,

  quantity      integer,
  shipping_date timestamptz,
  status        text CHECK (
    status IN ('placed', 'approved', 'delivered')
  ),
  complete      boolean DEFAULT FALSE,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.orders;
-- +goose StatementEnd
