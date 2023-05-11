package main

import (
	"context"
	"fmt"
	"github.com/lwenjim/study-golang/code/demo9/server-streaming/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type repoService struct {
	service.UnimplementedRepoServer
}

type userService struct {
	service.UnimplementedUsersServer
}

func (s userService) GetUser(
	ctx context.Context,
	in *service.UserGetRequest,
) (*service.UserGetReply, error) {
	log.Printf("Received request for user with Email: %s Id: %s\n", in.Email, in.Id)
	components := strings.Split(in.Email, "@")
	if len(components) != 2 {
		return nil, status.Error(codes.InvalidArgument, "Invalid email address specified")
	}
	u := service.User{
		Id:        in.Id,
		FirstName: components[0],
		LastName:  components[1],
		Age:       36,
	}
	return &service.UserGetReply{
		User: &u,
	}, nil
}

func (receiver repoService) GetRepos(
	in *service.RepoGetRequest,
	stream service.Repo_GetReposServer,
) error {
	log.Printf("Received request for repo with Created:%s Id:%s\n", in.CreatorId, in.Id)
	repo := service.Repository{
		Id: in.Id,
		Owner: &service.User{
			Id:        in.CreatorId,
			FirstName: "Jane",
		},
	}
	cnt := 1

	for {
		repo.Name = fmt.Sprintf("repo-%d", cnt)
		repo.Url = fmt.Sprintf("https://git.example.com/test/%s", repo.Name)

		r := service.RepoGetReply{
			Repo: &repo,
		}
		if err := stream.Send(&r); err != nil {
			return err
		}
		if cnt >= 5 {
			break
		}
		cnt++
	}
	return nil
}

func (receiver repoService) CreateBuild(
	in *service.Repository,
	stream service.Repo_CreateBuildServer,
) error {
	log.Printf(
		"Received build request for repo: %s\n", in.Name,
	)
	logLine := service.RepoBuildLog{
		LogLine:   fmt.Sprintf("Starting build for repository:%s", in.Name),
		Timestamp: timestamppb.Now(),
	}
	if err := stream.Send(&logLine); err != nil {
		return err
	}
	cnt := 1

	for {
		logLine := service.RepoBuildLog{
			LogLine:   fmt.Sprintf("Build log line  - %d", cnt),
			Timestamp: timestamppb.Now(),
		}
		if err := stream.Send(&logLine); err != nil {
			return err
		}
		if cnt >= 3 {
			break
		}
		time.Sleep(300 * time.Millisecond)
		cnt++
	}
	logLine = service.RepoBuildLog{
		LogLine:   fmt.Sprintf("Finished build for repository:%s", in.Name),
		Timestamp: timestamppb.Now(),
	}
	if err := stream.Send(&logLine); err != nil {
		return err
	}
	return nil
}

func registerServices(
	s *grpc.Server,
) {
	service.RegisterUsersServer(s, &userService{})
	service.RegisterRepoServer(s, &repoService{})
}

func startServer(
	s *grpc.Server,
	l net.Listener,
) error {
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
	registerServices(s)
	startServer(s, lis)
}
