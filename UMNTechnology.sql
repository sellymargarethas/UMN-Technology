CREATE DATABASE UMNTechnology;

-- DROP DATABASE IF EXISTS UMNTechnology;
USE DATABASE UMNTechnology;

-- DROP TABLE IF EXISTS users;
-- DROP SEQUENCE IF EXISTS users_seq;
CREATE SEQUENCE IF NOT EXISTS users_seq INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 999999999 CACHE 1;

CREATE TABLE
    users (
        id INTEGER PRIMARY KEY NOT NULL DEFAULT nextval ('users_seq'),
        nama VARCHAR(255) NOT NULL,
        email VARCHAR(100) NOT NULL,
        username VARCHAR(50) NOT NULL,
        password VARCHAR(100) NOT NULL,
        createdAt TIMESTAMP NOT NULL,
        updatedAt TIMESTAMP,
        deletedAt TIMESTAMP
    );

ALTER SEQUENCE users_seq OWNED BY users.id;