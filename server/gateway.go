package server

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/r1zq1/grpcsimplebank/pb"
	"google.golang.org/grpc"
)

func StartGatewayServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	err := pb.RegisterTransferServiceHandlerFromEndpoint(
		ctx, mux, "localhost:50051", []grpc.DialOption{grpc.WithInsecure()},
	)
	if err != nil {
		log.Fatalf("Gateway registration failed: %v", err)
	}

	log.Println("🌐 REST Gateway listening on :8080")
	http.ListenAndServe(":8080", mux)
}
