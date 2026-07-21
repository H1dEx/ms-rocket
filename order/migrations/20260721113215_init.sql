-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id serial primary key,
    order_uuid text not null,
    user_uuid text not null,
    part_uuids text[] not null,
    total_price real not null,
    transaction_uuid text,
    payment_method text,
    status text not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
