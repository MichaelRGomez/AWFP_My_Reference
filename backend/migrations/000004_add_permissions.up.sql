-- Filename: MyReference/backend/migrations/000004_add_permissions.up.sql
create table if not exists permissions(
    id bigserial primary key,
    code text not null
);

create table if not exists users_permissions(
    user_id bigint not null references users (id) on delete cascade,
    permissions_id bigint not null references permissions (id) on delete cascade,
    primary key (user_id, permissions_id)
);

insert into permissions (code)
values ('schools:read'), ('schools:write');