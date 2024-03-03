-- +goose Up

CREATE TABLE subscriber_channels (
    user_id UUID NOT NULL REFERENCES subscribers(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,

    PRIMARY KEY (user_id, title)
);

-- +goose Down
DROP TABLE subscriber_channels;