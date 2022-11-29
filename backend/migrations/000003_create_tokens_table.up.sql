-- Filename: MyReference/backend/migrations/000003_create_tokens_table.up.sql
create table if not exists tokens(
    hash bytea primary key,
    user_id bigint not null references users on delete cascade,
    expiry timestamp(0) with time zone not null,
    scope text not null
);