package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/lwenjim/study-golang/code/demo8/user-service/service1"
	"google.golang.org/grpc"
)

func setupGrpcConnection(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}

func getUserServiceClient(conn *grpc.ClientConn) service1.UsersClient {
	return service1.NewUsersClient(conn)
}

func getUser(client service1.UsersClient, u *service1.UserGetRequest) (*service1.UserGetReply, error) {
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

	result, err := getUser(c, &service1.UserGetRequest{
		Email: "jane@doe.com",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "User: %s %s\n", result.User.FirstName, result.User.LasttName)
}
