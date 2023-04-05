-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id INT GENERATED ALWAYS AS IDENTITY,
    name varchar(255),
    PRIMARY KEY(id)
);

INSERT INTO
    users (id, name) OVERRIDING USER VALUE
VALUES
    (1, 'admin');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
