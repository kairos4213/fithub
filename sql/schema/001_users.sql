-- +goose Up
create table users (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    first_name varchar(100) not null,
    middle_name varchar(100) default null,
    last_name varchar(100) not null,
    email varchar(254) unique not null,
    hashed_password text not null,
    profile_image varchar(255) default null,
    preferences json default null,
    is_admin boolean default false not null
);

-- +goose statementbegin
insert into users (
    id,
    created_at,
    updated_at,
    first_name,
    last_name,
    email,
    hashed_password,
    is_admin
) values (
    gen_random_uuid(),
    now(),
    now(),
    'user',
    'test',
    'user@email.com',
    '$argon2id$v=19$m=65536,t=1,p=6$9doHnQdcfE3W2945paPvbA$6OQd1ACMSsdyDY/p1ohZ0+WD6Hrl9WnfB7IVu/r4kjg',
    false
);
-- +goose statementend

-- +goose Down
drop table users;
