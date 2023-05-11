package cmd

import (
	"context"
	"github.com/lwenjim/study-golang/code/demo8/user-service/client-json/exercise1/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"strings"
	"testing"
)

type dummyUserService struct {
	service.UnimplementedUsersServer
}

func (s *dummyUserService) GetUser(ctx context.Context, in *service.UserGetRequest) (*service.UserGetReply, error) {
	components := strings.Split(in.Email, "@")
	u := service.User{
		Id:        in.Id,
		FirstName: components[0],
		LastName:  components[1],
		Age:       36,
	}
	return &service.UserGetReply{User: &u}, nil
}

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	service.RegisterUsersServer(s, &dummyUserService{})
	go func() {
		err := s.Serve(l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}

func TestCallUserSvc(t *testing.T) {
	testConfigs := []struct {
		c        grpConfig
		respJson string
		err      error
	}{
		{
			c:        grpConfig{method: "GetUser", request: `{"email":"john@doe.com","id":"user-123"}`},
			err:      nil,
			respJson: `{"user":{"id":"user-123","firstName":"john","lastName":"doe.com","age":36}}`,
		},
		{
			c:        grpConfig{},
			err:      ErrInvalidGrpcMethod,
			respJson: "",
		},
		{
			c:        grpConfig{method: "GetUser", request: "foo-bar"},
			err:      InvalidInputError{},
			respJson: "",
		},
	}

	s, l := startTestGrpcServer()
	defer s.GracefulStop()

	bufconnDialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return l.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(bufconnDialer))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range testConfigs {
		t.Log(tc)
		usersClient := getUserServiceClient(conn)
		respJson, err := callUsersMethod(usersClient, tc.c)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}
		if tc.err != nil && err == nil {
			t.Fatalf("Expected non-nil error, got nil")
		}
		//if tc.err != nil && !errors.As(err, &tc.err) {
		//	t.Fatalf("Expected error: %v, got: %v", tc.err, err)
		//}
		sanitizedRespJson := strings.Replace(string(respJson), " ", "", -1)
		sanitizedRespJson = strings.Replace(sanitizedRespJson, "\n", "", -1)

		if sanitizedRespJson != tc.respJson {
			t.Fatalf("Expected result: %v Got: %v", tc.respJson, sanitizedRespJson)
		}
	}

}
