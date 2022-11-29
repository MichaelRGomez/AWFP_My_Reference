-- Filename: MyReference/backend/migrations/000002_create_users_table.up.sql
create table if not exists users(
    id bigserial primary key,
    create_at timestamp(0) with time zone not null default now(),
    name text not null,
    email citext unique not null,
    password_hash bytea not null,
    activated bool not null,
    version integer not null default 1
);