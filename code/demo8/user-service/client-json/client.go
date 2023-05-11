package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
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

func createUserRequest(jsonQuery string) (*service.UserGetRequest, error) {
	u := service.UserGetRequest{}
	input := []byte(jsonQuery)
	return &u, protojson.Unmarshal(input, &u)
}

func getUserResponseJson(result *service.UserGetReply) ([]byte, error) {
	return protojson.Marshal(result)
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Must specify a gRPC server address and search query")
	}
	serverAddr := os.Args[1]
	u, err := createUserRequest(os.Args[2])
	if err != nil {
		log.Fatalf("Bad user input: %v", err)
	}

	conn, err := setupGrpcConnection(serverAddr)

	c := getUserServiceClient(conn)

	result, err := getUser(c, u)
	if err != nil {
		s := status.Convert(err)
		if s.Code() != codes.OK {
			log.Fatalf("Request failed: %v-%v\n", s.Code(), s.Message())
		}
		//log.Fatal(err)
	}
	data, err := getUserResponseJson(result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(os.Stdout, string(data))
}
