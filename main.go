package main

import (
	"log"
	"net"

	"github.com/andantan/vote-blockchain-server/network"
	"github.com/andantan/vote-blockchain-server/vote"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")

	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	s := network.Server{}

	grpcServer := grpc.NewServer()

	vote.RegisterBlockchainServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to server gRPC server over port 9000: %v", err)
	}

}
