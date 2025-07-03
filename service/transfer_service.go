package service

import (
	"context"
	"time"

	db "github.com/r1zq1/grpcsimplebank/db/sqlc"
	"github.com/r1zq1/grpcsimplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransferServer struct {
	pb.UnimplementedTransferServiceServer
	Store db.Store
}

func (s *TransferServer) Transfer(ctx context.Context,
	req *pb.TransferRequest) (*pb.TransferResponse, error) {
	arg := db.TrasferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Amount,
	}

	transfer, err := s.Store.TrasferTx(ctx, arg)
	if err != nil {
		return nil,
			status.Errorf(codes.Internal,
				"gagal membuat transfer: %v", err)
	}

	resp := &pb.TransferResponse{
		Id:            transfer.Transfer.ID,
		FromAccountId: transfer.Transfer.FromAccountID,
		ToAccountId:   transfer.Transfer.ToAccountID,
		Amount:        transfer.Transfer.Amount,
		CreatedAt:     transfer.Transfer.CreatedAt.Time.Format(time.RFC3339),
	}
	return resp, nil
}
