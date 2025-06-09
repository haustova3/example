package users

import (
	"context"
	"example/comments/internal/external/api/v1"
	mwc "example/comments/internal/external/middlewares"
	"example/comments/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserService struct {
	Client external.UsersClient
}

func NewUsersService(ctx context.Context, usersAddress string) (*UserService, error) {
	logger.Infow(ctx, "start users client", "address", usersAddress)
	conn, err := grpc.NewClient(usersAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			mwc.Logger,
			mwc.Tracer,
			mwc.Validate))
	if err != nil {
		logger.Errorw(ctx, "users service unavailable", "error", err.Error())
		return nil, err
	}
	usersClient := external.NewUsersClient(conn)
	return &UserService{
		Client: usersClient,
	}, nil
}

func (s *UserService) CheckUserID(ctx context.Context, userID int64) (bool, error) {
	req := &external.CheckUserIDRequest{
		UserID: userID,
	}
	res, err := s.Client.CheckUserID(ctx, req)
	if err != nil {
		return false, err
	}
	return res.IsCorrect, err
}
