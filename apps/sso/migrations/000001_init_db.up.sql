BEGIN;

create table if not exists users
(
    id              uuid                                   not null,
    first_name      text                                   not null,
    last_name       text                                   not null,
    email           text                                   not null,
    hashed_password text                                   not null,
    created_at      timestamp with time zone default now() not null,
    updated_at      timestamp with time zone default now() not null,
    constraint users_pk primary key (id),
    constraint users_email_key unique (email)
);

COMMIT;
