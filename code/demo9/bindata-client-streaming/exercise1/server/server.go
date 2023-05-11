package main

import (
	"fmt"
	"github.com/lwenjim/study-golang/code/demo9/bindata-client-streaming/exercise1/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"os"
)

type repoService struct {
	service.UnimplementedRepoServer
}

func (receiver repoService) CreateRepo(
	stream service.Repo_CreateRepoServer,
) error {
	var repoContext *service.RepoContext
	var data []byte

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(
				codes.Unknown,
				err.Error(),
			)
		}
		switch t := r.Body.(type) {
		case *service.RepoCreateRequest_Context:
			repoContext = r.GetContext()
		case *service.RepoCreateRequest_Data:
			d := r.GetData()
			data = append(data, d...)
		case nil:
			return status.Errorf(
				codes.InvalidArgument,
				"Message doesn't contain context or data",
			)
		default:
			return status.Errorf(
				codes.FailedPrecondition,
				"Unexpected message type:%s",
				t,
			)
		}
	}
	repo := service.Repository{
		Name: repoContext.Name,
		Url: fmt.Sprintf(
			"https://git.example.com/%s/%s",
			repoContext.CreatorId,
			repoContext.Name,
		),
	}
	r := service.RepoCreateReply{
		Repo: &repo,
		Size: int32(len(data)),
	}
	log.Printf("%#v\n", string(data))
	return stream.SendAndClose(&r)
}

func registerService(s *grpc.Server) {
	service.RegisterRepoServer(s, repoService{})
}

func startServer(s *grpc.Server, l net.Listener) error {
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

	s := grpc.NewServer()
	registerService(s)
	log.Fatal(startServer(s, lis))
}
