package main

import (
	"context"
	"github.com/lwenjim/study-golang/code/demo10/user-service-tls/service"
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

func startTestGrpcServer() (server *grpc.Server, listener *bufconn.Listener) {
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
	client, err := grpc.DialContext(
		context.Background(),
		"localhost",
		credsOption,
		grpc.WithContextDialer(bufconnDialer),
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
		t.Errorf(
			"Expected FirstName to be: jane, Got: %s",
			resp.User.FirstName,
		)
	}

}
