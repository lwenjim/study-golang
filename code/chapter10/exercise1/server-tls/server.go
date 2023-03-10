package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/lwenjim/code/chapter10/exercise1/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io"
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
		"Received request for user with Email: %s Id: %s\n",
		in.Email,
		in.Id,
	)

	components := strings.Split(in.Email, "@")
	if len(components) != 2 {
		return nil, errors.New("invalid email address")
	}
	u := service.UserGetReply{
		User: &service.User{
			Id:        in.Id,
			FirstName: components[0],
			LastName:  components[1],
			Age:       36,
		},
	}
	return &u, nil
}

func (receiver userService) GetHelp(
	stream service.Users_GetHelpServer,
) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf(
			"Request received: %s\n",
			request.Request,
		)
		response := service.UserHelpReply{
			Response: request.Request,
		}
		err = stream.Send(&response)
		if err != nil {
			return err
		}
	}
	return nil
}

func registerServices(
	s *grpc.Server,
	h *health.Server,
) {
	service.RegisterUsersServer(s, &userService{})
	grpc_health_v1.RegisterHealthServer(s, h)
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

func startServer(
	s *grpc.Server,
	l net.Listener,
) error {
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
	if len(tlsKeyFile) == 0 || len(tlsCertFile) == 0 {
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
	h := health.NewServer()
	registerServices(s, h)
	updateServiceHealth(
		h,
		service.Users_ServiceDesc.ServiceName,
		grpc_health_v1.HealthCheckResponse_SERVING,
	)
	log.Fatal(startServer(s, lis))
}
