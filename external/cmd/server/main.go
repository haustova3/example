package main

import (
	"example/external/internal/server"
	"example/external/internal/server/mw"
	desc "example/external/pkg/api/v1"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	log.Println("App starting")
	if err := listenAndServe(); err != nil {
		log.Fatalf("App error: %e", err)
	}
}

func listenAndServe() error {

	address := fmt.Sprintf(":%s", "8093")
	list, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen and serve app failed: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.Logger,
			mw.Validate,
		),
	)

	reflection.Register(grpcServer)
	controller := &server.Controller{}
	desc.RegisterUsersServer(grpcServer, controller)
	desc.RegisterProductsServer(grpcServer, controller)
	if err = grpcServer.Serve(list); err != nil {
		log.Fatalf("Server err: %e", err)
	}
	return nil
}
