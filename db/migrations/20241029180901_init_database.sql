-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
 id UUID NOT NULL PRIMARY KEY,
 username VARCHAR(255) NOT NULL,
 email varchar(255) NOT NULL,
 hashed_password varchar(255) NOT NULL
);
INSERT INTO users(id, username, email, hashed_password) VALUES ('388f37f1-1a5a-4439-9c41-155f4e2470cf', 'username', 'email', 'hash');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
