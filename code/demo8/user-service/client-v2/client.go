package main

import (
	"context"
	"fmt"
	"github.com/lwenjim/study-golang/code/demo8/user-service/service"
	"google.golang.org/grpc"
	"log"
	"os"
)

func setupGrpcConnection(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}

func getUserServiceClient(conn *grpc.ClientConn) service.UsersClient {
	return service.NewUsersClient(conn)
}

func getUser(client service.UsersClient, u *service.UserGetRequest) (*service.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Must specify a gRPC server address")
	}
	conn, err := setupGrpcConnection(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := getUserServiceClient(conn)

	result, err := getUser(c, &service.UserGetRequest{
		Email: "jane@doe.com",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "%#v: %s\n", result.User, result.Location)
}
