package service

import (
	"context"
	"fmt"

	"github.com/r1zq1/grpcsimplebank/config"
	"github.com/r1zq1/grpcsimplebank/pb"
	"github.com/r1zq1/grpcsimplebank/worker"

	db "github.com/r1zq1/grpcsimplebank/db/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountServer struct {
	pb.UnimplementedAccountServiceServer
	Store           db.Store
	TaskDistributor worker.TaskDistributor
	Config          config.Config
}

func NewAccountServer(store db.Store, distributor worker.TaskDistributor, config config.Config) *AccountServer {
	return &AccountServer{
		Store:           store,
		TaskDistributor: distributor,
		Config:          config,
	}
}

func (s *AccountServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	// Simpan ke DB
	arg := db.CreateAccountParams{
		Owner:   req.Owner,
		Balance: req.Balance,
		Email:   req.Email,
	}

	account, err := s.Store.CreateAccount(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create account: %v", err)
	}

	// Kirim task async ke Redis untuk email
	taskPayload := worker.PayloadSendEmail{
		Email: req.Email,
		Owner: req.Owner,
	}

	err = s.TaskDistributor.DistributeSendWelcomeEmail(ctx, taskPayload)
	if err != nil {
		// Kirim task gagal tetap response success, log saja
		fmt.Printf("failed to send welcome email: %v\n", err)
	}

	resp := &pb.CreateAccountResponse{
		Id:      account.ID,
		Owner:   account.Owner,
		Email:   account.Email,
		Balance: account.Balance,
	}

	return resp, nil
}
