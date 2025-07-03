package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/r1zq1/grpcsimplebank/config"
	db "github.com/r1zq1/grpcsimplebank/db/sqlc"
	"github.com/r1zq1/grpcsimplebank/pb"
	"github.com/r1zq1/grpcsimplebank/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("gagal baca config: %v", err)
	}

	conn, err := sql.Open("postgres", conf.DBSource)
	if err != nil {
		log.Fatalf("tidak bisa connect ke db: %v", err)
	}

	store := db.NewStore(conn)
	server := grpc.NewServer()

	pb.RegisterTransferServiceServer(server,
		&service.TransferServer{Store: store})
	reflection.Register(server)

	listener, err := net.Listen("tcp", conf.GRPCAddress)
	if err != nil {
		log.Fatalf("gagal listen: %v", err)
	}

	log.Printf("server berjalan di %s", conf.GRPCAddress)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("gagal serve: %v", err)
	}
}
