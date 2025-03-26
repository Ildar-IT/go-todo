-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, email, password_hash) VALUES ('Test','test@test.ru', '6d7973616c74b1b3773a05c0ed0176787a4f1574ff0075f7521e');

INSERT INTO user_roles (user_id, role_id) VALUES ((SELECT id FROM users WHERE email = 'test@test.ru'), 1)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DELETE FROM users WHERE email = 'test@test.ru'
-- +goose StatementEnd
