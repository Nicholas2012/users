-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION pg_trgm;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    surname VARCHAR NOT NULL,
    patronymic VARCHAR NOT NULL DEFAULT '',
    age INT,
    gender CHAR(1),
    nationality VARCHAR
);
CREATE INDEX users_names ON users USING gin ((name || ' ' || surname || ' ' || patronymic) gin_trgm_ops);
CREATE INDEX users_age ON users (age);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
