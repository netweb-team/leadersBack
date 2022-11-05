create table users (
    id serial primary key,
    login text not null unique,
    pass bytea not null
);

create table tables (
    id serial primary key,
    path text not null unique
);

create table patterns (
    pool_id int,
    pattern int,
    lng real not null,
    lat real not null,
    avg_price real not null default 0,
    primary key(pool_id, pattern)
);

create table analogs (
    id serial primary key,
    lng real not null,
    lat real not null,
    addr text not null,
    room text not null,
    segment text not null,
    floors int not null,
    cur_floor int not null,
    walls text not null,
    total real not null,
    kitchen real not null,
    balcony text not null,
    metro real not null,
    state text not null,
    price int not null,
    avg_price real not null,
    sale_coef real not null,
    floor_coef real not null,
    total_coef real not null,
    kitchen_coef real not null,
    balcony_coef real not null,
    metro_coef real not null,
    state_coef real not null,
    pool int not null,
    pattern int not null,
    use boolean not null default 't'
);

create table archive (
    id serial primary key,
    pool int not null,
    pattern int not null,
    price real not null,
    analogs json,
    coefs json,
    price_path text not null
);
