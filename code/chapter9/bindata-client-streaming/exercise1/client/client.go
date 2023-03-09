package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/lwenjim/code/chapter9/bindata-client-streaming/exercise1/service"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
)

type appConfig struct {
	filePath string
	sererUrl string
}

func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}

func getRepoServiceClient(conn *grpc.ClientConn) service.RepoClient {
	return service.NewRepoClient(conn)
}

func setupFlags(w io.Writer, args []string) (appConfig, error) {
	c := appConfig{}
	fs := flag.NewFlagSet("grpc-client", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.filePath, "file-path", "", "Repository contents to upload")
	err := fs.Parse(args)
	if err != nil {
		return c, err
	}
	if len(c.filePath) == 0 {
		return c, errors.New("file path empty")
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return c, errors.New("must specify server URL as the only positional argument")
	}
	c.sererUrl = fs.Arg(0)

	return c, nil
}

func uploadRepository(r io.Reader, repoClient service.RepoClient) (*service.RepoCreateReply, error) {
	stream, err := repoClient.CreateRepo(
		context.Background(),
	)
	if err != nil {
		return nil, err
	}
	contextData := service.RepoCreateRequest_Context{
		Context: &service.RepoContext{
			CreatorId: "user-123",
			Name:      "test-repo",
		},
	}
	request := service.RepoCreateRequest{
		Body: &contextData,
	}
	err = stream.Send(&request)
	if err != nil {
		return nil, err
	}
	size := 32 * 1024
	buf := make([]byte, size)
	for {

		nBytes, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		bData := service.RepoCreateRequest_Data{
			Data: buf[:nBytes],
		}
		request = service.RepoCreateRequest{
			Body: &bData,
		}
		err = stream.Send(&request)
		if err != nil {
			return nil, err
		}
	}
	return stream.CloseAndRecv()
}

func main() {
	c, err := setupFlags(os.Stdout, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	conn, err := setupGrpcConn(c.sererUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	f, err := os.Open(c.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	repoClient := getRepoServiceClient(conn)
	resp, err := uploadRepository(f, repoClient)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"Uploaded %d bytes\n",
		resp.Size,
	)
}
