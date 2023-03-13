package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/lwenjim/code/chapter10/exercise2/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var shutdownTimeout time.Duration = 30 * time.Second

type userService struct {
	service.UnimplementedUsersServer
}

func (s userService) GetUser(
	ctx context.Context,
	in *service.UserGetRequest,
) (*service.UserGetReply, error) {
	log.Printf("Received request for user with Email: %s Id: %s\n", in.Email, in.Id)
	components := strings.Split(in.Email, "@")
	if len(components) != 2 {
		return nil, errors.New("invalid email address")
	}
	if components[0] == "panic" {
		panic("I was asked to panic")
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

func registerServices(s *grpc.Server, h *health.Server) {
	service.RegisterUsersServer(s, &userService{})
	grpc_health_v1.RegisterHealthServer(s, h)
}

func startServer(s *grpc.Server, l net.Listener) error {
	return s.Serve(l)
}

func updateServiceHealth(
	h *health.Server,
	service string,
	status grpc_health_v1.HealthCheckResponse_ServingStatus,
) {
	h.SetServingStatus(
		service,
		status,
	)
}

func shutDown(
	s *grpc.Server, waitForShutdownCompletion chan struct{},
) {
	s.GracefulStop()
	waitForShutdownCompletion <- struct{}{}
	log.Println("shutdown")
}

func waitForShutDown(
	ctx context.Context,
	s *grpc.Server,
	h *health.Server,
) {
	waitForShutdownCompletion := make(chan struct{})
	sctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-sctx.Done()

	log.Printf(
		"Got signal: %v. Server shutting down.",
		ctx,
	)

	updateServiceHealth(
		h,
		service.Users_ServiceDesc.ServiceName,
		grpc_health_v1.HealthCheckResponse_NOT_SERVING,
	)

	go shutDown(s, waitForShutdownCompletion)

	select {
	case <-ctx.Done():
		log.Println("Forcing shutdown")
		s.Stop()
	case <-waitForShutdownCompletion:
		return
	}
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
	h := health.NewServer()
	registerServices(s, h)
	updateServiceHealth(
		h,
		service.Users_ServiceDesc.ServiceName,
		grpc_health_v1.HealthCheckResponse_SERVING,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()
	go waitForShutDown(ctx, s, h)
	log.Fatalf(
		"Shutting down: %v\n",
		startServer(s, lis),
	)
}
