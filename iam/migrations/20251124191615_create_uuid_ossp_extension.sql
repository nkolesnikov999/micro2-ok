-- +goose Up
-- создаем расширение uuid, если оно еще не установлено
create extension if not exists "uuid-ossp";

-- +goose Down
-- удаляем расширение uuid, если оно не используется
drop extension if exists "uuid-ossp";
