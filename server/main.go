package main

import (
	cont "grpc-server/controller/account"
	pb "grpc-server/proto/account"
	"log"
	"net"

	"google.golang.org/grpc"
)

const port = ":50051"

func main() {

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatal("Failed to listen")
	}

	s := grpc.NewServer()

	pb.RegisterDepositServiceServer(s, &cont.DepositServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}
