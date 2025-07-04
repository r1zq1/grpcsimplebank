package main

import (
	"fmt"
	"log"

	"github.com/r1zq1/grpcsimplebank/client"
	"github.com/r1zq1/grpcsimplebank/pb"
)

func main() {
	// transfer unary
	grpcClient, err := client.NewGRPCClient("localhost:50051")
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer grpcClient.Close()

	res, err := grpcClient.Transfer(1, 2, 500)
	if err != nil {
		log.Fatalf("transfer failed: %v", err)
	}

	fmt.Printf("✅ Transfer ID: %d, From: %d, To: %d, Amount: %d, Time: %s\n",
		res.Id, res.FromAccountId, res.ToAccountId, res.Amount, res.CreatedAt)

	// transfer_history – Server Streaming Example
	cli, err := client.NewGRPCClient("localhost:50051")
	if err != nil {
		log.Fatal("cannot connect:", err)
	}
	defer cli.Close()

	if err := cli.GetTransferHistory(1); err != nil {
		log.Fatalf("stream error: %v", err)
	}

	// batch_transfer – Client Streaming Example
	cli3, err3 := client.NewGRPCClient("localhost:50051")
	if err3 != nil {
		log.Fatal("cannot connect:", err3)
	}
	defer cli3.Close()

	transfers3 := []pb.TransferRequest{
		{FromAccountId: 1, ToAccountId: 2, Amount: 100},
		{FromAccountId: 1, ToAccountId: 3, Amount: 200},
		{FromAccountId: 1, ToAccountId: 4, Amount: 300},
	}

	if err := cli.BatchTransfer(transfers3); err != nil {
		log.Fatalf("batch error: %v", err)
	}

	// live_transfer – Bidirectional Streaming Example
	cli4, err4 := client.NewGRPCClient("localhost:50051")
	if err != nil {
		log.Fatal("cannot connect:", err4)
	}
	defer cli4.Close()

	transfers4 := []pb.TransferRequest{
		{FromAccountId: 2, ToAccountId: 1, Amount: 50},
		{FromAccountId: 2, ToAccountId: 3, Amount: 150},
		{FromAccountId: 2, ToAccountId: 4, Amount: 250},
	}

	if err := cli.LiveTransfer(transfers4); err != nil {
		log.Fatalf("live error: %v", err)
	}
}
