-- Filename: MyReference/backend/migrations/000005_create_reference_info_table.up.sql

create table if not exists reference_info(
  id bigserial primary key,
  create_at timestamp(0) with time zone not null default now(),
  name text not null,
  location text not null,
  version integer not null default 1  
);