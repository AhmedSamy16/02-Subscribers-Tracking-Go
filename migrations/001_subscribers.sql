-- +goose Up

CREATE TABLE subscribers (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);
-- +goose Down
DROP TABLE subscribers;