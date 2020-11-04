-- +goose Up
CREATE TABLE todo (
    id UUID PRIMARY KEY,
    title varchar(20) NOT NULL,
    content text NOT NULL
);

-- +goose Down
DROP TABLE todo;
