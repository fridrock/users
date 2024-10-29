-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
 id UUID NOT NULL PRIMARY KEY,
 username VARCHAR(255) NOT NULL,
 email varchar(255) NOT NULL,
 -- TODO make them NOT NULL
 name VARCHAR(255),
 surname VARCHAR(255),
 hashed_password varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS friends(
    fr1id UUID NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    fr2id UUID NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT friends_pkey PRIMARY KEY (fr1id, fr2id)
);
INSERT INTO users(id, username, email, hashed_password) 
VALUES ('388f37f1-1a5a-4439-9c41-155f4e2470cf', 'user1', 'email1', 'password1'), 
('6c5d4660-d821-4e55-ab2f-03cccdc1b446', 'user2', 'email2', 'password2');
INSERT INTO friends(fr1id, fr2id) VALUES ('388f37f1-1a5a-4439-9c41-155f4e2470cf', '6c5d4660-d821-4e55-ab2f-03cccdc1b446');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS friends;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
