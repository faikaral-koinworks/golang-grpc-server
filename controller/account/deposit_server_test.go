package account_test

import (
	"context"
	"log"
	"net"
	"testing"

	cont "grpc-server/controller/account"
	pb "grpc-server/proto/account"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterDepositServiceServer(server, &cont.DepositServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

}

func TestDepositServiceServer_GetDeposit(t *testing.T) {
	test := struct {
		name string
		res  *pb.GetDepositResponse
		err  error
	}{
		"Testing getting total deposit",
		&pb.GetDepositResponse{TotalDeposit: 0},
		nil,
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewDepositServiceClient(conn)

	t.Run(test.name, func(t *testing.T) {
		request := &pb.GetDepositRequest{}

		response, err := client.GetDeposit(ctx, request)

		if response.TotalDeposit != test.res.TotalDeposit {
			t.Error("response expected", test.res, "received", response)
		}

		if err != nil {
			t.Error("response expected", test.err, "received", err)
		}
	})

}

func TestDepositServiceServer_Deposit(t *testing.T) {
	test := []struct {
		name   string
		amount float32
		res    *pb.DepositResponse
		errMsg string
	}{
		{
			"Invalid Request with negative amount",
			-1.1,
			&pb.DepositResponse{Ok: false},
			"cannot deposit negative balance",
		},
		{
			"Valid Request",
			1.1,
			&pb.DepositResponse{Ok: true},
			"",
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewDepositServiceClient(conn)

	for _, v := range test {
		t.Run(v.name, func(t *testing.T) {
			request := &pb.DepositRequest{Amount: v.amount}

			response, err := client.Deposit(ctx, request)

			if response.GetOk() != v.res.GetOk() {
				t.Error("response expected", v.res.GetOk(), "received", response.GetOk())
			}

			if err != nil {
				if er, _ := status.FromError(err); er.Message() != v.errMsg {
					t.Error("error code:expected", v.errMsg, "received", err)
				}
			}
		})
	}

}
