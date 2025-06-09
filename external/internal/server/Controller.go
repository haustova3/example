package server

import (
	"context"
	servicepb "example/external/pkg/api/v1"
)

type Controller struct {
	servicepb.UnimplementedUsersServer
	servicepb.UnimplementedProductsServer
}

func (s *Controller) CheckUserID(_ context.Context, in *servicepb.CheckUserIDRequest) (*servicepb.CheckUserIDResponse, error) {
	return &servicepb.CheckUserIDResponse{
		IsCorrect: in.UserID < 100,
	}, nil
}

func (s *Controller) GetOwner(_ context.Context, in *servicepb.GetOwnerRequest) (*servicepb.GetOwnerResponse, error) {
	if in.ProductID > 100 {
		return &servicepb.GetOwnerResponse{
			OwnerID: 0,
		}, nil
	}
	return &servicepb.GetOwnerResponse{
		OwnerID: in.ProductID + 73,
	}, nil
}
