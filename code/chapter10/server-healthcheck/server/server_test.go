package main

import (
	"context"
	"github.com/lwenjim/code/chapter10/server-healthcheck/service"
	"google.golang.org/grpc"
	"net"
	"testing"
)

func TestUserService(t *testing.T) {
	l := startTestGrpcServer()
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
			Email: "jane@doe.com",
			Id:    "foo-bar",
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
