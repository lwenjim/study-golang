package main

import (
	"context"
	"github.com/lwenjim/code/chapter10/svc-client-dxn/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	registerServices(s)
	go func() {
		err := startServer(s, l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}
func TestUserService(t *testing.T) {
	s, l := startTestGrpcServer()
	defer s.GracefulStop()

	bufConnDialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufConnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}
	usersClient := service.NewUsersClient(client)
	resp, err := usersClient.GetUser(
		context.Background(),
		&service.UserGetRequest{
			Email: "jane@doe.com",
			Id:    "foo-bar",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.User.FirstName != "jane" {
		t.Errorf("Expected FirstName to be: jane, Got: %s", resp.User.FirstName)
	}
}
