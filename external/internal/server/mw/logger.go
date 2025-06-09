package mw

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	raw, _ := protojson.Marshal((req).(proto.Message))
	log.Printf("grpc request received method: %s, req: %s", info.FullMethod, string(raw))
	if resp, err = handler(ctx, req); err != nil {
		log.Printf("grpc response error. Method: %s, rerr: %s", info.FullMethod, err.Error())
		return
	}
	rawResp, _ := protojson.Marshal((resp).(proto.Message))
	log.Printf("grpc response. Method: %s, resp: %s", info.FullMethod, string(rawResp))
	return
}
