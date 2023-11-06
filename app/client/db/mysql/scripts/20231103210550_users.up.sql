CREATE TABLE IF NOT EXISTS users (
    id BIGINT NOT NULL AUTO_INCREMENT,
    name varchar(64) NOT NULL,
    email varchar(128) NOT NULL UNIQUE,
    password varchar(256) NOT NULL,
    PRIMARY KEY(id)
);
