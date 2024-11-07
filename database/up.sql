DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id varchar(32) PRIMARY KEY,
    password varchar(255) not null,
    email varchar(255) not null,
    created_at timestamp not null default NOW()
);

DROP TABLE IF EXISTS posts;

CREATE TABLE posts (
    id varchar(32) PRIMARY KEY,
    post_content varchar(32) not null,
    created_at timestamp not null default NOW(),
    user_id varchar(32) not null,
    FOREIGN KEY(user_id) REFERENCES users(id)
);
