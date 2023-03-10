package main

import (
	"context"
	"errors"
	"github.com/lwenjim/code/chapter10/user-service-tls/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"strings"
)

type userService struct {
	service.UnimplementedUsersServer
}

func (receiver userService) GetUser(
	ctx context.Context,
	in *service.UserGetRequest,
) (*service.UserGetReply, error) {
	log.Printf(
		"Received request for user with Email: %s Id:%s\n",
		in.Email,
		in.Id,
	)
	components := strings.Split(in.Email, "@")
	if len(components) != 2 {
		return nil, errors.New("invalid email address")
	}
	u := service.User{

		Id:        in.Id,
		FirstName: components[0],
		LastName:  components[1],
		Age:       36,
	}
	return &service.UserGetReply{
		User: &u,
	}, nil
}

func registerServices(s *grpc.Server) {
	service.RegisterUsersServer(s, &userService{})
}

func startServer(s *grpc.Server, l net.Listener) error {
	return s.Serve(l)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8880"
	}
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	tlsCertFile := os.Getenv("TLS_CERT_FILE_PATH")
	tlsKeyFile := os.Getenv("TLS_KEY_FILE_PATH")

	if len(tlsCertFile) == 0 || len(tlsKeyFile) == 0 {
		log.Fatal("TLS_CERT_FILE_PATH and TLS_KEY_FILE_PATH must both be specified")
	}

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
	log.Fatal(startServer(s, lis))
}
