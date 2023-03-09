package main

import (
	"fmt"
	"github.com/lwenjim/code/chapter9/bidi-streaming/service"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
)

type userService struct {
	service.UnimplementedUsersServer
}

func (receiver userService) GetHep(
	stream service.Users_GetHepServer,
) error {
	log.Println("Client connected")
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("Request received: %s \n", request.Request)

		response := service.UserHelpReply{
			Response: request.Request,
		}
		err = stream.Send(&response)
		if err != nil {
			return err
		}
	}
	log.Println("Client Disconnected")
	return nil
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
	s := grpc.NewServer()
	registerServices(s)
	log.Fatal(startServer(s, lis))
}
