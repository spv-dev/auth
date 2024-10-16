-- +goose Up
-- +goose StatementBegin
create table users (
    id bigserial primary key,
    name varchar(255) not null,
    email varchar(255) not null,
    password varchar(255) not null,
    role smallint not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd