package main

import (
	"context"
	"fmt"
	"github.com/lwenjim/study-golang/code/demo8/user-service/service2"
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

func getUserServiceClient(conn *grpc.ClientConn) service2.UsersClient {
	return service2.NewUsersClient(conn)
}

func getUser(client service2.UsersClient, u *service2.UserGetRequest) (*service2.UserGetReply, error) {
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

	result, err := getUser(c, &service2.UserGetRequest{
		Email: "jane@doe.com",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "%#v: %s\n", result.User, result.Location)
}
