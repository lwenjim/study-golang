package main

import (
	"context"
	"fmt"
	"github.com/lwenjim/study-golang/code/demo8/user-service/service1"
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
func getUserServiceClient(conn *grpc.ClientConn) service1.UsersClient {
	return service1.NewUsersClient(conn)
}

func getUser(client service1.UsersClient, u *service1.UserGetRequest) (*service1.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func createUserRequest(jsonQuery string) (*service1.UserGetRequest, error) {
	u := service1.UserGetRequest{}
	input := []byte(jsonQuery)
	return &u, protojson.Unmarshal(input, &u)
}

func getUserResponseJson(result *service1.UserGetReply) ([]byte, error) {
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
