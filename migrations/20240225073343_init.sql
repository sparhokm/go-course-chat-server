-- +goose Up
create table chats
(
    id         serial primary key,
    user_ids integer[] not null,
    created_at timestamp      not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table chats;