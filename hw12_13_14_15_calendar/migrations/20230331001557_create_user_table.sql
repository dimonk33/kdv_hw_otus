-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id INT GENERATED ALWAYS AS IDENTITY,
    name varchar(255),
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
