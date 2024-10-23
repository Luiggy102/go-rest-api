DROP TABLE IF EXISTS user;

CREATE TABLE user (
    id varchar(32) PRIMARY KEY,
    password varchar(255) not null,
    email varchar(255) not null,
    created_at timestamp not null default NOW(),
);
