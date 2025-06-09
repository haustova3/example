package usecases

import (
	"context"
	"errors"
	"example/comments/internal/model"
)

type SaveCommentRepository interface {
	SaveComment(_ context.Context, comment model.Comment) (int64, error)
}

type UserService interface {
	CheckUserID(_ context.Context, userID int64) (bool, error)
}

type ProductsService interface {
	GetProductOwner(_ context.Context, productID int64) (int64, error)
}

type CreateCommentService struct {
	rep            SaveCommentRepository
	userService    UserService
	productService ProductsService
}

func NewCreateCommentService(rep SaveCommentRepository, products ProductsService, users UserService) *CreateCommentService {
	return &CreateCommentService{
		rep:            rep,
		productService: products,
		userService:    users,
	}
}

func (s *CreateCommentService) CreateComment(ctx context.Context, comment model.Comment) (int64, error) {
	isCorrectUserID, err := s.userService.CheckUserID(ctx, comment.UserID)
	if err != nil {
		return 0, errors.Join(model.ErrUserServiceUnavailable, err)
	}
	if !isCorrectUserID {
		return 0, model.ErrIncorrectUserID
	}
	productOwnerID, err := s.productService.GetProductOwner(ctx, comment.ProductID)
	if err != nil {
		return 0, errors.Join(model.ErrProductServiceUnavailable, err)
	}
	if productOwnerID == 0 {
		return 0, model.ErrProductOwnerNotFound
	}
	comment.ProductOwnerID = productOwnerID
	commentID, err := s.rep.SaveComment(ctx, comment)
	return commentID, err
}
