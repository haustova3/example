package usecases

import (
	"context"
	"example/comments/internal/model"
)

type GetCommentsRepository interface {
	GetComments(_ context.Context, productID int64) ([]model.Comment, error)
}

type GetCommentsService struct {
	rep GetCommentsRepository
}

func NewGetCommentsService(rep GetCommentsRepository) *GetCommentsService {
	return &GetCommentsService{
		rep: rep,
	}
}

func (service *GetCommentsService) GetComments(ctx context.Context, productID int64) ([]model.Comment, error) {
	commentID, err := service.rep.GetComments(ctx, productID)
	return commentID, err
}
