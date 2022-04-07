package account

import (
	"context"
	"errors"
	pb "grpc-server/proto/account"
	"log"
)

var accountValue float32 = 0

type DepositServer struct {
	pb.UnimplementedDepositServiceServer
}

func (d *DepositServer) Deposit(c context.Context, in *pb.DepositRequest) (*pb.DepositResponse, error) {
	log.Printf("Received deposit request with the amount of: %v", in.GetAmount())
	if in.GetAmount() <= 0 {
		return &pb.DepositResponse{Ok: false}, errors.New("cannot deposit negative balance")
	}
	accountValue = accountValue + in.GetAmount()
	return &pb.DepositResponse{Ok: true}, nil
}

func (d *DepositServer) GetDeposit(c context.Context, in *pb.GetDepositRequest) (*pb.GetDepositResponse, error) {
	log.Printf("Received RequestGetDeposit")
	return &pb.GetDepositResponse{TotalDeposit: accountValue}, nil
}
