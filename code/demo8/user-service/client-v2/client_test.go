package main

import (
	"context"
	"github.com/lwenjim/study-golang/code/demo8/user-service/service2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

type dummyUserService struct {
	service2.UnimplementedUsersServer
}

func (s *dummyUserService) GetUser(ctx context.Context, in *service2.UserGetRequest) (*service2.UserGetReply, error) {
	u := service2.User{
		Id:        "user-123-a",
		FirstName: "jane",
		LasttName: "doe",
		Age:       36,
	}
	return &service2.UserGetReply{
		User: &u,
	}, nil
}

func startServer(s *grpc.Server, l net.Listener) error {
	return s.Serve(l)
}

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	service2.RegisterUsersServer(s, &dummyUserService{})
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

	result, err := c.GetUser(context.Background(), &service2.UserGetRequest{Email: "jane@doe.com"})

	if err != nil {
		t.Fatal(err)
	}

	if result.User.FirstName != "jane" || result.User.LasttName != "doe" {
		t.Fatalf("Expected: jane doe, Got: %s %s", result.User.FirstName, result.User.LasttName)
	}
}
