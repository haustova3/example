package mw

import (
	"context"
	"example/comments/internal/logger"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	raw, _ := protojson.Marshal((req).(proto.Message))
	logger.Infow(ctx, "grpc request received", "method", info.FullMethod, "req", string(raw))
	if resp, err = handler(ctx, req); err != nil {
		logger.Warnw(ctx, "grpc err response", "method", info.FullMethod, "err", err.Error())
		return
	}
	rawResp, _ := protojson.Marshal((resp).(proto.Message))
	logger.Infow(ctx, "grpc response", "method", info.FullMethod, "resp", string(rawResp))
	return
}
