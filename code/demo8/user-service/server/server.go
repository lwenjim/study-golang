package main

import (
	"context"
	"errors"
	"github.com/lwenjim/study-golang/code/demo8/user-service/service2"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strings"
)

type userService struct {
	service2.UnimplementedUsersServer
}

func (s userService) GetUser(ctx context.Context, in *service2.UserGetRequest) (*service2.UserGetReply, error) {
	log.Printf("Received request for user with Email: %s Id: %s\n", in.Email, in.Id)
	components := strings.Split(in.Email, "@")
	if len(components) != 2 {
		return nil, errors.New("invalid email address")
	}
	u := service2.User{
		Id:        in.Id,
		FirstName: components[0],
		LasttName: components[1],
		Age:       36,
	}
	return &service2.UserGetReply{
		User:     &u,
		Location: "abc",
	}, nil
}

func registerServices(s *grpc.Server) {
	service2.RegisterUsersServer(s, &userService{})
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
	s := grpc.NewServer()
	registerServices(s)
	log.Fatal(startServer(s, lis))
}
