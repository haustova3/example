package mwc

import (
	"context"
	"example/comments/internal/trace"

	"google.golang.org/grpc"
)

func Tracer(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, _ ...grpc.CallOption) error {
	ctxWithSpan, span := trace.Tracer().
		Start(
			ctx,
			method,
		)
	defer span.End()
	return invoker(ctxWithSpan, method, req, reply, cc)
}
