package main

import (
	"context"
	"fmt"
	"github.com/lwenjim/study-golang/code/demo9/server-streaming/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"log"
	"net"
	"testing"
)

func startTestGrpcServer() (
	*grpc.Server,
	*bufconn.Listener,
) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	registerServices(s)
	go func() {
		err := startServer(s, l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}

func TestUserService(t *testing.T) {
	_, l := startTestGrpcServer()
	bufconnDialer := func(
		ctx context.Context,
		addr string,
	) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(
		context.Background(),
		"", grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}

	usersClient := service.NewUsersClient(client)
	resp, err := usersClient.GetUser(
		context.Background(),
		&service.UserGetRequest{
			Id:    "foo-bar",
			Email: "jane@doe.com",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.User.FirstName != "jane" {
		t.Errorf(
			"Expected FirstName to be: jane, Got: %s",
			resp.User.FirstName,
		)
	}

}

func TestRepoService(t *testing.T) {
	_, l := startTestGrpcServer()
	bufconnDialer := func(
		ctx context.Context,
		addr string,
	) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}
	repoClient := service.NewRepoClient(client)
	stream, err := repoClient.GetRepos(context.Background(), &service.RepoGetRequest{
		CreatorId: "user-123",
		Id:        "repo-123",
	})
	if err != nil {
		t.Fatal(err)
	}
	var repos []*service.Repository
	for {
		repo, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		repos = append(repos, repo.Repo)
	}
	if len(repos) != 5 {
		t.Fatalf("Expected to get back 5 repos, got back:%d repos", len(repos))
	}

	for idx, repo := range repos {
		gotRepoName := repo.Name
		expectedRepoName := fmt.Sprintf("repo-%d", idx+1)
		if gotRepoName != expectedRepoName {
			t.Errorf("Expected Repo Name to be:%s, Got:%s", expectedRepoName, gotRepoName)
		}
	}
}

func TestRepoBuildMethod(t *testing.T) {
	_, l := startTestGrpcServer()

	bufconnDialer := func(
		ctx context.Context,
		addr string,
	) (net.Conn, error) {
		return l.Dial()
	}
	client, err := grpc.DialContext(
		context.Background(),
		"", grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}

	repoClient := service.NewRepoClient(client)
	stream, err := repoClient.CreateBuild(
		context.Background(),
		&service.Repository{
			Name: "practicalgo/test/repo",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	var logLines []*service.RepoBuildLog
	for {
		line, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		logLines = append(logLines, line)
	}
	if len(logLines) != 5 {
		t.Fatalf(
			"Expected to get back 3 lines in the log, got back:%d repos",
			len(logLines),
		)
	}
	expectedLastLine := "Starting build for repository:practicalgo/test/repo"
	if logLines[0].LogLine != expectedLastLine {
		t.Fatalf(
			"Expected first line to be:%s, Got:%s",
			expectedLastLine,
			logLines[0].LogLine,
		)
	}
	expectedLastLine = "Finished build for repository:practicalgo/test/repo"
	if logLines[4].LogLine != expectedLastLine {
		t.Fatalf(
			"Expected last line to be:%s,Got:%s",
			expectedLastLine,
			logLines[4].LogLine,
		)
	}

	logLine := logLines[0]
	if err := logLine.Timestamp.CheckValid(); err != nil {
		t.Fatalf(
			"Logline timestamp invalid: %#v",
			logLine,
		)
	}
}
