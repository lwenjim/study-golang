package main

import (
	"context"
	"fmt"
	"github.com/lwenjim/code/chapter10/client-resiliency/service"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"time"
)

func setupGrpcConn(addr string) (context.CancelFunc, *grpc.ClientConn, error) {
	log.Printf("Connecting to server on %s\n", addr)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithReturnConnectionError(),
	)
	return cancel, conn, err
}
func getUserServiceClient(conn *grpc.ClientConn) service.UsersClient {
	return service.NewUsersClient(conn)
}

func getUser(
	client service.UsersClient,
	u *service.UserGetRequest,
) (
	*service.UserGetReply,
	context.CancelFunc,
	error,
) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	resp, err := client.GetUser(ctx, u, grpc.WaitForReady(true))
	return resp, cancel, err
}

func setupChat(
	r io.Reader,
	w io.Writer,
	c service.UsersClient,
) (err error) {
	var clientConn = make(chan service.Users_GetHelpClient)
	var done = make(chan bool)

	stream, err := createHelpStream(c)
	defer stream.CloseSend()

	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	go func() {
		for {
			clientConn <- stream
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
			}

			if err != nil {
				log.Printf("Recreating stream.")
			}
			stream, err = createHelpStream(c)
			if err != nil {
				close(clientConn)
				done <- true
			} else {
				fmt.Printf(
					"Response:%s\n",
					resp.Response,
				)
				if resp.Response == "hello-10" {
					done <- true
				}
			}
		}
	}()
	requestMsg := "hello"
	msgCount := 1
	for {
		if msgCount > 10 {
			break
		}
		stream = <-clientConn
		if stream == nil {
			break
		}
		request := service.UserHelpRequest{
			Request: fmt.Sprintf("%s-%d", requestMsg, msgCount),
		}
		err := stream.Send(&request)
		if err != nil {
			log.Printf("Send error:%v, Will retry.\n", err)
		} else {
			log.Printf("Request sent: %d\n", msgCount)
			msgCount++
		}
	}

	<-done
	return stream.CloseSend()
}

func createHelpStream(c service.UsersClient) (service.Users_GetHelpClient, error) {
	return c.GetHelp(
		context.Background(),
		grpc.WaitForReady(true),
	)
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Specify a gRPC server and method to call")
	}
	serverAddr := os.Args[1]
	methodName := os.Args[2]

	cancel, conn, err := setupGrpcConn(serverAddr)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := getUserServiceClient(conn)

	switch methodName {
	case "GetUser":
		for i := 0; i < 5; i++ {
			log.Printf("Request: %d\n", i)
			userEmail := os.Args[3]
			result, cancel, err := getUser(c, &service.UserGetRequest{
				Email: userEmail,
			})
			defer cancel()
			if err != nil {
				log.Fatalf("getUser failed: %v", err)
			}
			fmt.Fprintf(
				os.Stdout,
				"User: %s %s\n",
				result.User.FirstName,
				result.User.LastName,
			)
			time.Sleep(1 * time.Second)
		}
	case "GetHelp":
		err = setupChat(os.Stdin, os.Stdout, c)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Unrecognized method name")
	}
}
