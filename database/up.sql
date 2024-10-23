DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id varchar(32) PRIMARY KEY,
    password varchar(255) not null,
    email varchar(255) not null,
    created_at timestamp not null default NOW()
);
