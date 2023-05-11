package main

import (
	"context"
	"net"
	"testing"

	"github.com/lwenjim/study-golang/code/demo8/user-service/service1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type dummyUserService struct {
	service1.UnimplementedUsersServer
}

func (s *dummyUserService) GetUser(ctx context.Context, in *service1.UserGetRequest) (*service1.UserGetReply, error) {
	u := service1.User{
		Id:        "user-123-a",
		FirstName: "jane",
		LasttName: "doe",
		Age:       36,
	}
	return &service1.UserGetReply{
		User: &u,
	}, nil
}

func startServer(s *grpc.Server, l net.Listener) error {
	return s.Serve(l)
}

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	service1.RegisterUsersServer(s, &dummyUserService{})
	go func() {
		startServer(s, l)
	}()
	return s, l
}

func TestGetUser(t *testing.T) {
	s, l := startTestGrpcServer()
	defer s.GracefulStop()
	bufConnDialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return l.Dial()
	}
	conn, err := grpc.DialContext(
		context.Background(),
		"", grpc.WithInsecure(),
		grpc.WithContextDialer(bufConnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}

	c := getUserServiceClient(conn)

	result, err := c.GetUser(context.Background(), &service1.UserGetRequest{Email: "jane@doe.com"})

	if err != nil {
		t.Fatal(err)
	}

	if result.User.FirstName != "jane" || result.User.LasttName != "doe" {
		t.Fatalf("Expected: jane doe, Got: %s %s", result.User.FirstName, result.User.LasttName)
	}
}
