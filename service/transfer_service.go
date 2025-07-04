package service

import (
	"context"
	"io"
	"log"
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

// Unary Transfer
func (s *TransferServer) Transfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {

	log.Printf("Transfer from %d to %d amount %d", req.FromAccountId, req.ToAccountId, req.Amount)

	if req.FromAccountId == 0 || req.ToAccountId == 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			"from_account and to_account must not be empty")
	}

	if req.Amount <= 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			"amount must be positive")
	}

	if req.FromAccountId == req.ToAccountId {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot transfer to the same account")
	}

	// Simulasi rekening tidak ditemukan
	if req.FromAccountId == -1 {
		return nil, status.Errorf(codes.NotFound,
			"source account not found")
	}

	result, err := s.Store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Amount,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer failed: %v", err)
	}

	return &pb.TransferResponse{
		Id:            result.Transfer.ID,
		FromAccountId: result.Transfer.FromAccountID,
		ToAccountId:   result.Transfer.ToAccountID,
		Amount:        result.Transfer.Amount,
		CreatedAt:     result.Transfer.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}

// Server-Streaming
func (s *TransferServer) GetTransferHistory(req *pb.HistoryRequest, stream pb.TransferService_GetTransferHistoryServer) error {
	transfers, err := s.Store.ListTransfersByAccountID(stream.Context(), req.AccountId)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get history: %v", err)
	}

	for _, tr := range transfers {
		res := &pb.TransferResponse{
			Id:            tr.ID,
			FromAccountId: tr.FromAccountID,
			ToAccountId:   tr.ToAccountID,
			Amount:        tr.Amount,
			CreatedAt:     tr.CreatedAt.Time.Format(time.RFC3339),
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}

	return nil
}

// Client-Streaming
func (s *TransferServer) BatchTransfer(stream pb.TransferService_BatchTransferServer) error {
	var success, failed int32

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.TransferSummary{SuccessCount: success, FailedCount: failed})
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "receive error: %v", err)
		}

		_, err = s.Store.TransferTx(stream.Context(), db.TransferTxParams{
			FromAccountID: req.FromAccountId,
			ToAccountID:   req.ToAccountId,
			Amount:        req.Amount,
		})

		if err != nil {
			failed++
		} else {
			success++
		}
	}
}

// Bi-Directional Streaming
func (s *TransferServer) LiveTransfer(stream pb.TransferService_LiveTransferServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "recv failed: %v", err)
		}

		result, err := s.Store.TransferTx(stream.Context(), db.TransferTxParams{
			FromAccountID: req.FromAccountId,
			ToAccountID:   req.ToAccountId,
			Amount:        req.Amount,
		})

		res := &pb.TransferResponse{}
		if err == nil {
			res = &pb.TransferResponse{
				Id:            result.Transfer.ID,
				FromAccountId: result.Transfer.FromAccountID,
				ToAccountId:   result.Transfer.ToAccountID,
				Amount:        result.Transfer.Amount,
				CreatedAt:     result.Transfer.CreatedAt.Time.Format(time.RFC3339),
			}
		}

		if err := stream.Send(res); err != nil {
			return err
		}
	}
}

// func (s *TransferServer) Transfer(ctx context.Context,
// 	req *pb.TransferRequest) (*pb.TransferResponse, error) {
// 	arg := db.TransferTxParams{
// 		FromAccountID: req.FromAccountId,
// 		ToAccountID:   req.ToAccountId,
// 		Amount:        req.Amount,
// 	}

// 	transfer, err := s.Store.TransferTx(ctx, arg)
// 	if err != nil {
// 		return nil,
// 			status.Errorf(codes.Internal,
// 				"gagal membuat transfer: %v", err)
// 	}

// 	resp := &pb.TransferResponse{
// 		Id:            transfer.Transfer.ID,
// 		FromAccountId: transfer.Transfer.FromAccountID,
// 		ToAccountId:   transfer.Transfer.ToAccountID,
// 		Amount:        transfer.Transfer.Amount,
// 		CreatedAt:     transfer.Transfer.CreatedAt.Time.Format(time.RFC3339),
// 	}
// 	return resp, nil
// }
