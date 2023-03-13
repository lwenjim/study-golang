package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/lwenjim/code/chapter10/client-resiliency/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type userService struct {
	service.UnimplementedUsersServer
}

func (s userService) GetUser(
	ctx context.Context,
	in *service.UserGetRequest,
) (*service.UserGetReply, error) {
	if ctx.Err() != nil {
		log.Printf(
			"Request context canceled.Returning.",
		)
		return nil, status.Errorf(
			codes.Canceled,
			"Request Canceled",
		)
	}

	log.Printf(
		"Received request for user with Email: %s Id: %s\n",
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
	log.Printf("Returning Response.")
	return &service.UserGetReply{User: &u}, nil
}

func (s userService) GetHelp(
	stream service.Users_GetHelpServer,
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
		if request.Request == "panic" {
			panic("I was asked to panic")
		}
		response := service.UserHelpReply{
			Response: request.Request,
		}
		err = stream.Send(&response)
		if err != nil {
			return err
		}
	}
	log.Println("Client disconnected")
	return nil
}

func logMessage(
	ctx context.Context,
	method string,
	latency time.Duration,
	err error,
) {
	var requestId string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Print("No metadata")
	} else {
		if len(md.Get("Request-Id")) != 0 {
			requestId = md.Get("Request-Id")[0]
		}
	}
	log.Printf(
		"Method:%s, Duration:%s, Error:%v, Request-Id:%s",
		method,
		latency,
		err,
		requestId,
	)
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
	s := grpc.NewServer(
	//grpc.ChainUnaryInterceptor(
	//	clientDisconnectUnaryInterceptor,
	//),
	//grpc.ChainStreamInterceptor(
	//	clientDisconnectStreamInterceptor,
	//),
	)
	registerServices(s)
	log.Fatal(startServer(s, lis))
}
