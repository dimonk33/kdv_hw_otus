-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
                       id int NOT NULL,

                       title varchar(255),
                       start_time   timestamp,
                       end_time     timestamp,
                       description text,
                       own_user_id   int,

                       PRIMARY KEY(id),
                       CONSTRAINT fk_user
                           FOREIGN KEY(own_user_id)
                               REFERENCES hw.public.users(id)

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP CONSTRAINT fk_user;
DROP TABLE events;
-- +goose StatementEnd
