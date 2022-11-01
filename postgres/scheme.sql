create table users (
    id serial primary key,
    login text not null unique,
    pass bytea not null
);

create table tables (
    id serial primary key,
    path text not null unique
);