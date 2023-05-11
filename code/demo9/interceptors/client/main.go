package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/lwenjim/study-golang/code/demo9/interceptors/service"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
)

func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(
			metadataUnaryInterceptor,
			loggingUnaryInterceptor,
		),
		grpc.WithChainStreamInterceptor(
			metadataStreamInterceptor,
			loggingStreamingInterceptor,
		),
	)
}
func getUserServiceClient(conn *grpc.ClientConn) service.UsersClient {
	return service.NewUsersClient(conn)
}
func getUser(
	client service.UsersClient,
	u *service.UserGetRequest,
) (
	*service.UserGetReply,
	error,
) {
	return client.GetUser(context.Background(), u)
}
func setupChat(
	r io.Reader,
	w io.Writer,
	c service.UsersClient,
) error {
	stream, err := c.GetHelp(context.Background())
	if err != nil {
		return err
	}
	for {
		scanner := bufio.NewScanner(r)
		prompt := "Request: "
		fmt.Fprint(w, prompt)

		scanner.Scan()

		if err := scanner.Err(); err != nil {
			return err
		}
		msg := scanner.Text()
		if msg == "quit" {
			break
		}
		request := service.UserHelpRequest{
			Request: msg,
		}
		err := stream.Send(&request)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("Response: %s\n", resp.Response)
	}
	return stream.CloseSend()
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Specify a gRPC server and method to call")
	}
	serverAddr := os.Args[1]
	methodName := os.Args[2]

	conn, err := setupGrpcConn(serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := getUserServiceClient(conn)

	switch methodName {
	case "GetUser":
		result, err := getUser(c, &service.UserGetRequest{
			Email: "jane@doe.com",
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(
			os.Stdout,
			"User: %s %s\n",
			result.User.FirstName,
			result.User.LasttName,
		)
	case "GetHelp":
		err = setupChat(os.Stdin, os.Stdout, c)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Unrecognized method name")
	}
}
