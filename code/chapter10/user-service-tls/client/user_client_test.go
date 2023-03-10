package main

import (
	"context"
	"github.com/lwenjim/code/chapter10/user-service-tls/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

var (
	tlsCertFile = "../tls/server.crt"
	tlsKeyFile  = "../tls/server.key"
)

type dummyUserService struct {
	service.UnimplementedUsersServer
}

func (receiver dummyUserService) GetUser(
	ctx context.Context,
	in *service.UserGetRequest,
) (
	*service.UserGetReply,
	error,
) {
	u := service.User{
		Id:        "user-123-a",
		FirstName: "jane",
		LastName:  "doe",
		Age:       36,
	}
	return &service.UserGetReply{
		User: &u,
	}, nil
}

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(20)
	creds, err := credentials.NewServerTLSFromFile(
		tlsCertFile,
		tlsKeyFile,
	)
	if err != nil {
		log.Fatal(err)
	}
	credsOption := grpc.Creds(creds)
	s := grpc.NewServer(credsOption)
	service.RegisterUsersServer(s, &dummyUserService{})

	go func() {
		err := s.Serve(l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}

func TestGetUser(t *testing.T) {
	s, l := startTestGrpcServer()
	defer s.GracefulStop()

	bufconnDialer := func(
		ctx context.Context,
		addr string,
	) (net.Conn, error) {
		return l.Dial()
	}

	creds, err := credentials.NewClientTLSFromFile(
		tlsCertFile,
		"",
	)
	if err != nil {
		t.Fatal(err)
	}

	credsOption := grpc.WithTransportCredentials(creds)
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost",
		credsOption,
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}

	c := getUserServiceClient(conn)
	result, err := getUser(
		c,
		&service.UserGetRequest{Email: "jane@doe.com"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if result.User.FirstName != "jane" || result.User.LastName != "doe" {
		t.Fatalf(
			"Expected: jane doe, Got: %s %s",
			result.User.FirstName,
			result.User.LastName,
		)
	}
}
