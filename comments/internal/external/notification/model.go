package notification

import "time"

type CommentNotification struct {
	ID        int64     `json:"id"`
	OwnerID   int64     `json:"owner_id"`
	CommentID int64     `json:"comment_id"`
	CreatedTS time.Time `json:"operation_time"`
}
