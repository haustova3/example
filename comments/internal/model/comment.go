package model

import "time"

type Comment struct {
	ID             int64
	ProductID      int64
	ProductOwnerID int64
	UserID         int64
	Text           string
	Ts             time.Time
}
