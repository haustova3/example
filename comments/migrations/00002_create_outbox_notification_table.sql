-- +goose Up
-- +goose StatementBegin
CREATE TABLE outbox_notification
(
    id         bigserial PRIMARY KEY,
    owner_id   bigint    not null,
    comment_id bigint    not null,
    ts         timestamp not null DEFAULT now(),
    status     text      not null DEFAULT 'new'
);
ALTER TABLE outbox_notification
    ADD CONSTRAINT product_id_positive CHECK ( owner_id > 0 );
ALTER TABLE outbox_notification
    ADD CONSTRAINT user_id_positive CHECK ( comment_id > 0 );
ALTER TABLE outbox_notification
    ADD CONSTRAINT check_status CHECK ( status IN ('new', 'send'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE outbox_notification;
-- +goose StatementEnd

