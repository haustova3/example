package repository

import (
	"context"
	"example/comments/internal/external/notification"
	"example/comments/internal/model"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	write    *pgxpool.Pool
	ntfCount int32
}

func NewRepository(write *pgxpool.Pool, ntfCount int) *Repository {
	return &Repository{write: write,
		ntfCount: int32(ntfCount)}
}

func (rep *Repository) SaveComment(ctx context.Context, comment model.Comment) (int64, error) {
	var err error
	commentID := int64(0)
	err = pgx.BeginTxFunc(ctx, rep.write, pgx.TxOptions{}, func(tx pgx.Tx) error {
		createdTS := time.Now()
		r := New(tx)
		commentID, err = r.SaveComment(ctx, &SaveCommentParams{
			UserID:    comment.UserID,
			ProductID: comment.ProductID,
			Tx:        comment.Text,
			Ts: pgtype.Timestamp{
				Time:  createdTS,
				Valid: true,
			},
		})
		if err != nil {
			return fmt.Errorf("save comment faild: %w", err)
		}
		err = rep.saveNotification(ctx, r, comment.ProductOwnerID, commentID, createdTS)
		if err != nil {
			return fmt.Errorf("comment create ntf failed: %w", err)
		}
		return nil
	})
	return commentID, err
}

func (rep *Repository) GetComments(ctx context.Context, productID int64) ([]model.Comment, error) {
	r := New(rep.write)
	comments, err := r.GetCommentsByProduct(ctx, productID)
	if err != nil {
		return make([]model.Comment, 0), err
	}
	res := make([]model.Comment, len(comments))
	for i, val := range comments {
		res[i] = model.Comment{
			ID:        val.ID,
			UserID:    val.UserID,
			ProductID: productID,
			Text:      val.Tx,
			Ts:        val.Ts.Time,
		}
	}
	return res, nil
}

func (rep *Repository) saveNotification(ctx context.Context, r *Queries, ownerID int64, commentID int64, createdTS time.Time) error {
	err := r.SaveNotification(ctx, &SaveNotificationParams{
		OwnerID:   ownerID,
		CommentID: commentID,
		Ts: pgtype.Timestamp{
			Time:  createdTS,
			Valid: true,
		},
	})
	return err
}

func (rep *Repository) GetCommentNotification(ctx context.Context) ([]notification.CommentNotification, error) {
	r := New(rep.write)
	ntfsEntity, err := r.GetUnSendNotification(ctx, rep.ntfCount)
	if err != nil {
		return []notification.CommentNotification{}, fmt.Errorf("can not get notification: %e", err)
	}
	ntfs := make([]notification.CommentNotification, len(ntfsEntity))
	for i, val := range ntfsEntity {
		ntfs[i] = notification.CommentNotification{
			ID:        val.ID,
			OwnerID:   val.OwnerID,
			CommentID: val.CommentID,
			CreatedTS: val.Ts.Time,
		}
	}
	return ntfs, nil
}
func (rep *Repository) MarkNotificationAsSend(ctx context.Context, notificationID int64) error {
	r := New(rep.write)
	err := r.MaskNotificationAsSend(ctx, notificationID)
	return err
}
