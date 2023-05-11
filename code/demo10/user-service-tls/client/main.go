package main

import (
	"context"
	"fmt"
	"github.com/lwenjim/study-golang/code/demo10/user-service-tls/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
)

func setupGrpcConn(
	addr string,
	tlsCertFile string,
) (
	*grpc.ClientConn,
	error,
) {
	creds, err := credentials.NewClientTLSFromFile(
		tlsCertFile,
		"",
	)
	if err != nil {
		return nil, err
	}
	credsOption := grpc.WithTransportCredentials(creds)
	return grpc.DialContext(
		context.Background(),
		addr,
		credsOption,
		grpc.WithBlock(),
	)
}

func getUserServiceClient(conn *grpc.ClientConn) service.UsersClient {
	return service.NewUsersClient(conn)
}

func getUser(
	client service.UsersClient,
	u *service.UserGetRequest,
) (*service.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(
			"Must specify a gRPC server address",
		)
	}

	tlsCertFile := os.Getenv("TLS_CERT_FILE_PATH")
	if len(tlsCertFile) == 0 {
		log.Fatal("TLS_CERT_FILE_PATH must be specified")
	}
	conn, err := setupGrpcConn(os.Args[1], tlsCertFile)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := getUserServiceClient(conn)

	result, err := getUser(
		c,
		&service.UserGetRequest{Email: "jane@doe.com"},
	)
	if err != nil {
		log.Printf("Error, getUser")
		log.Fatal(err)
	}
	fmt.Fprintf(
		os.Stdout,
		"User: %s %s\n",
		result.User.FirstName,
		result.User.LastName,
	)
}
