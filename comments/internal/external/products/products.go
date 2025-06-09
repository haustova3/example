package products

import (
	"context"
	"example/comments/internal/external/api/v1"
	mwc "example/comments/internal/external/middlewares"
	"example/comments/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductService struct {
	Client external.ProductsClient
}

func NewProductsService(ctx context.Context, productsAddress string) (*ProductService, error) {
	logger.Infow(ctx, "start products client", "address", productsAddress)
	conn, err := grpc.NewClient(productsAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			mwc.Logger,
			mwc.Tracer,
			mwc.Validate))
	if err != nil {
		logger.Errorw(ctx, "products service unavailable", "error", err.Error())
		return nil, err
	}
	productsClient := external.NewProductsClient(conn)
	return &ProductService{
		Client: productsClient,
	}, nil
}

func (s *ProductService) GetProductOwner(ctx context.Context, productID int64) (int64, error) {
	req := &external.GetOwnerRequest{
		ProductID: productID,
	}
	res, err := s.Client.GetOwner(ctx, req)
	if err != nil {
		return 0, err
	}
	return res.OwnerID, nil
}
