-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    uuid uuid primary key,
    login text not null unique,
    email text not null unique,
    password text not null,
    notification_methods jsonb not null default '[]',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
