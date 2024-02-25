-- +goose Up
create table chats
(
    id         serial primary key,
    user_names varchar(255)[] not null,
    created_at timestamp      not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table chats;