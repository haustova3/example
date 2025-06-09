package app

import (
	"context"
	"errors"
	"example/comments/internal/logger"
	"example/comments/internal/model"
	servicepb "example/comments/pkg/api/comments/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ servicepb.CommentsServer = (*CommentsController)(nil)

type CreateCommentService interface {
	CreateComment(ctx context.Context, comment model.Comment) (int64, error)
}

type GetCommentsService interface {
	GetComments(ctx context.Context, productID int64) ([]model.Comment, error)
}

type CommentsController struct {
	servicepb.UnimplementedCommentsServer
	createCommentService CreateCommentService
	getCommentsService   GetCommentsService
}

func NewCommentsController(createCommentService CreateCommentService,
	getCommentsService GetCommentsService,
) *CommentsController {

	return &CommentsController{
		createCommentService: createCommentService,
		getCommentsService:   getCommentsService,
	}
}

func (s *CommentsController) CreateComment(ctx context.Context, in *servicepb.CreateCommentRequest) (*servicepb.CreateCommentResponse, error) {
	comment := model.Comment{
		UserID:    in.UserID,
		ProductID: in.ProductID,
		Text:      in.Text,
	}
	commentID, err := s.createCommentService.CreateComment(ctx, comment)
	if err != nil {
		logger.Warnw(ctx, "Request failed", "error", err)
		if errors.Is(err, model.ErrIncorrectUserID) || errors.Is(err, model.ErrProductOwnerNotFound) {
			return nil, status.Error(codes.FailedPrecondition, "Invalid request")
		}
		if errors.Is(err, model.ErrProductServiceUnavailable) || errors.Is(err, model.ErrUserServiceUnavailable) {
			return nil, status.Error(codes.Unavailable, "External service unavailable")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}
	res := &servicepb.CreateCommentResponse{
		CommentID: commentID,
	}
	return res, nil
}

func (s *CommentsController) GetComments(ctx context.Context, in *servicepb.GetCommentsRequest) (*servicepb.GetCommentsResponse, error) {
	comments, err := s.getCommentsService.GetComments(ctx, in.ProductID)
	if err != nil {
		logger.Warnw(ctx, "Request failed", "error", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}
	commentsResponse := make([]*servicepb.Comment, len(comments))
	for i, val := range comments {
		commentsResponse[i] = &servicepb.Comment{
			ID:     val.ID,
			UserID: val.UserID,
			Text:   val.Text,
			Ts:     timestamppb.New(val.Ts),
		}
	}
	res := &servicepb.GetCommentsResponse{
		ProductID: in.ProductID,
		Comments:  commentsResponse,
	}
	return res, nil
}
