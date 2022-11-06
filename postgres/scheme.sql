create table users (
    id serial primary key,
    login text not null unique,
    pass bytea not null
);

create table cookies (
    user_id int,
    cookie text,
    primary key(user_id, cookie) 
);

create table tables (
    id serial primary key,
    path text not null unique,
    user_id int not null
);

create table patterns (
    pool_id int,
    pattern int,
    lng float8 not null,
    lat float8 not null,
    avg_price float8 not null default 0,
    primary key(pool_id, pattern)
);

create table analogs (
    id serial primary key,
    lng float8 not null,
    lat float8 not null,
    addr text not null,
    room text not null,
    segment text not null,
    floors int not null,
    cur_floor int not null,
    walls text not null,
    total float8 not null,
    kitchen float8 not null,
    balcony text not null,
    metro float8 not null,
    state text not null,
    price int not null,
    avg_price float8 not null,
    sale_coef float8 not null,
    floor_coef float8 not null,
    total_coef float8 not null,
    kitchen_coef float8 not null,
    balcony_coef float8 not null,
    metro_coef float8 not null,
    state_coef float8 not null,
    pool int not null,
    pattern int not null,
    use boolean not null default 't'
);

create table archive (
    id serial primary key,
    pool int not null,
    pattern int[] not null default '{}',
    price float8[] not null default '{}',
    analogs json[] not null default '{}',
    coefs json[] not null default '{}',
    price_path text not null,
    user_id int not null
);
