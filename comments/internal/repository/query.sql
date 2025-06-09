-- name: SaveComment :one
INSERT INTO comments (user_id, product_id, tx, ts)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetCommentsByProduct :many
SELECT id, user_id, tx, ts
FROM comments
WHERE product_id = $1;


-- name: SaveNotification :exec
INSERT INTO outbox_notification (owner_id, comment_id, ts)
VALUES ($1, $2, $3);

-- name: GetUnSendNotification :many
SELECT id, owner_id, comment_id, ts, status
FROM outbox_notification
WHERE status = 'new'
ORDER BY ts
    LIMIT $1;

-- name: MaskNotificationAsSend :exec
UPDATE outbox_notification
SET status = 'send'
WHERE id = $1;