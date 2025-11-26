-- +goose Up

-- создаем таблицу пользователей
create table if not exists users
(
    uuid                uuid primary key                  default uuid_generate_v4(),
    login               varchar(255)              not null unique,
    email               varchar(255)              not null unique,
    password_hash        varchar(255)              not null,
    notification_methods jsonb                     not null default '[]'::jsonb,
    created_at           timestamp with time zone not null default now(),
    updated_at           timestamp with time zone not null default now()
);

-- создаем индекс по login для быстрого поиска пользователя по логину
create index if not exists idx_users_login on users (login);

-- создаем индекс по email для быстрого поиска пользователя по email
create index if not exists idx_users_email on users (email);

-- создаем индекс по notification_methods для быстрого поиска по методам уведомлений (GIN индекс для JSONB)
create index if not exists idx_users_notification_methods on users using gin (notification_methods);

-- +goose Down
-- удаляем индексы
drop index if exists idx_users_notification_methods;
drop index if exists idx_users_email;
drop index if exists idx_users_login;

-- удаляем таблицу пользователей
drop table if exists users;