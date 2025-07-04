package server

import (
	"database/sql"
	"log"
	"net"

	"github.com/hibiken/asynq"
	"github.com/r1zq1/grpcsimplebank/config"
	db "github.com/r1zq1/grpcsimplebank/db/sqlc"
	"github.com/r1zq1/grpcsimplebank/pb"
	"github.com/r1zq1/grpcsimplebank/service"
	"github.com/r1zq1/grpcsimplebank/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGRPCServer() {
	conf, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("gagal baca config: %v", err)
	}

	conn, err := sql.Open("postgres", conf.DBSource)
	if err != nil {
		log.Fatalf("tidak bisa connect ke db: %v", err)
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: "localhost:6379",
	}

	store := db.NewStore(conn)
	server := grpc.NewServer()
	distributor := worker.NewRedisTaskDistributor(redisOpt)

	pb.RegisterTransferServiceServer(server,
		&service.TransferServer{Store: store})
	pb.RegisterAccountServiceServer(server,
		&service.AccountServer{
			Store:           store,
			TaskDistributor: distributor,
		})
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
