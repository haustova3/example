-- +goose Up
-- +goose StatementBegin
CREATE TABLE comments
(
    id         bigserial PRIMARY KEY,
    user_id    bigint    not null,
    product_id bigint    not null,
    tx         text      not null,
    ts         timestamp not null DEFAULT now()
);
CREATE INDEX product_id_idx ON comments USING HASH(product_id);
ALTER TABLE comments
    ADD CONSTRAINT user_id_positive CHECK ( user_id > 0 );
ALTER TABLE comments
    ADD CONSTRAINT product_id_positive CHECK ( product_id > 0 );
ALTER TABLE comments
    ADD CONSTRAINT text_len CHECK (LENGTH(tx) >= 5 AND LENGTH(tx) <= 255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE comments;
-- +goose StatementEnd
