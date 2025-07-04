package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/r1zq1/grpcsimplebank/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client pb.TransferServiceClient
}

func NewGRPCClient(address string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	client := pb.NewTransferServiceClient(conn)
	return &GRPCClient{conn: conn, client: client}, nil
}

func (c *GRPCClient) Close() {
	c.conn.Close()
}

// Unary transfer
func (c *GRPCClient) Transfer(fromID, toID, amount int64) (*pb.TransferResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.TransferRequest{
		FromAccountId: fromID,
		ToAccountId:   toID,
		Amount:        amount,
	}

	res, err := c.client.Transfer(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// --- Server Streaming ---
func (c *GRPCClient) GetTransferHistory(accountID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := c.client.GetTransferHistory(ctx, &pb.HistoryRequest{AccountId: accountID})
	if err != nil {
		return fmt.Errorf("error calling GetTransferHistory: %w", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("✅ History stream finished")
			break
		}
		if err != nil {
			return fmt.Errorf("stream recv error: %w", err)
		}

		log.Printf("📜 History: TransferID=%d | From=%d -> To=%d | Amount=%d | Time=%s\n",
			resp.Id, resp.FromAccountId, resp.ToAccountId, resp.Amount, resp.CreatedAt)
	}
	return nil
}

// --- Client Streaming ---
func (c *GRPCClient) BatchTransfer(transfers []pb.TransferRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := c.client.BatchTransfer(ctx)
	if err != nil {
		return fmt.Errorf("error opening stream: %w", err)
	}

	for _, tr := range transfers {
		if err := stream.Send(&tr); err != nil {
			return fmt.Errorf("send error: %w", err)
		}
		log.Printf("📤 Sent transfer from: %v to: %v amount: %v\n",
			tr.FromAccountId, tr.ToAccountId, tr.Amount)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("recv error: %w", err)
	}

	log.Printf("✅ Batch Summary: Success=%d, Failed=%d\n", resp.SuccessCount, resp.FailedCount)
	return nil
}

// --- Bidirectional Streaming ---
func (c *GRPCClient) LiveTransfer(transfers []pb.TransferRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stream, err := c.client.LiveTransfer(ctx)
	if err != nil {
		return fmt.Errorf("failed to open LiveTransfer stream: %w", err)
	}

	done := make(chan struct{})

	// Goroutine to receive responses
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("recv error: %v", err)
				break
			}
			log.Printf("📥 Transfer result: %+v", resp)
		}
		close(done)
	}()

	// Send all transfer requests
	for _, tr := range transfers {
		if err := stream.Send(&tr); err != nil {
			return fmt.Errorf("send error: %w", err)
		}
		log.Printf("📤 Sent: from: %v to: %v amount: %v\n",
			tr.FromAccountId, tr.ToAccountId, tr.Amount)
		time.Sleep(500 * time.Millisecond) // Simulate delay
	}

	// Close send side
	if err := stream.CloseSend(); err != nil {
		return fmt.Errorf("close send error: %w", err)
	}

	<-done
	log.Println("✅ Live stream completed")
	return nil
}
