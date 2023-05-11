package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/lwenjim/study-golang/code/demo8/multiple-services/server/exercise1/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type grpConfig struct {
	server      string
	method      string
	body        string
	request     string
	service     string
	prettyPrint bool
}

func validateGrpcConfig(c grpConfig) error {
	if len(c.service) == 0 {
		return errors.New("unrecognized")
	}
	if len(c.method) == 0 {
		return ErrInvalidGrpcMethod
	}
	return nil
}

func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
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

func createUserRequest(jsonQuery string) (*service.UserGetRequest, error) {
	u := service.UserGetRequest{}
	input := []byte(jsonQuery)
	return &u, protojson.Unmarshal(input, &u)
}

func getUser(client service.UsersClient, u *service.UserGetRequest) (*service.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func getUserResponseJson(c grpConfig, result *service.UserGetReply) ([]byte, error) {
	if c.prettyPrint {
		return []byte(protojson.Format(result)), nil
	}
	return protojson.Marshal(result)
}

func getRepoServiceClient(conn *grpc.ClientConn) service.RepoClient {
	return service.NewRepoClient(conn)
}

func callUsersMethod(usersClient service.UsersClient, c grpConfig) ([]byte, error) {

	switch c.method {
	case "GetUser":
		req, err := createUserRequest(c.request)
		if err != nil {
			return nil, InvalidInputError{Err: err}
		}
		result, err := getUser(usersClient, req)
		if err != nil {
			return nil, err
		}
		respData, err := getUserResponseJson(c, result)
		return respData, err
	case "":
		return nil, ErrInvalidGrpcMethod
	default:
		return nil, ErrInvalidGrpcMethod
	}
}

func createGetRepoRequest(jsonQuery string) (*service.RepoGetRequest, error) {
	v := service.RepoGetRequest{}
	input := []byte(jsonQuery)
	return &v, protojson.Unmarshal(input, &v)
}

func getRepos(client service.RepoClient, r *service.RepoGetRequest) (*service.RepoGetReply, error) {
	return client.GetRepos(context.Background(), r)
}

func getReposResponseJson(c grpConfig, result *service.RepoGetReply) ([]byte, error) {
	if c.prettyPrint {
		return []byte(protojson.Format(result)), nil
	}
	return protojson.Marshal(result)
}

func callRepoMethod(repoClient service.RepoClient, c grpConfig) ([]byte, error) {
	switch c.method {
	case "GetRepos":
		req, err := createGetRepoRequest(c.request)
		if err != nil {
			return nil, InvalidInputError{Err: err}
		}
		result, err := getRepos(repoClient, req)
		if err != nil {
			return nil, err
		}
		respData, err := getReposResponseJson(c, result)
		return respData, err
	case "":
		return nil, ErrInvalidGrpcMethod
	default:
		return nil, ErrInvalidGrpcMethod
	}
}

func HandleGrpc(w io.Writer, args []string) error {
	var err error
	c := grpConfig{}
	fs := flag.NewFlagSet("grpc", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.method, "method", "", "Method to call")
	fs.StringVar(&c.request, "request", "", "Request to send")
	fs.StringVar(&c.service, "service", "", "gRPC service to send the request to")
	fs.BoolVar(&c.prettyPrint, "pretty-print", false, "Pretty print the JSON output")
	fs.Usage = func() {
		var usageString = `
grpc: A gRPC client.

grpc: <options> server`
		fmt.Fprint(w, usageString)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}

	err = fs.Parse(args)
	if err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}
	c.server = fs.Arg(0)
	err = validateGrpcConfig(c)
	if err != nil {
		return err
	}

	conn, err := setupGrpcConn(c.server)
	if err != nil {
		return err
	}

	var usersClient service.UsersClient
	switch c.service {
	case "Users":
		usersClient = getUserServiceClient(conn)
		respJson, err := callUsersMethod(usersClient, c)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, string(respJson))
	case "Repo":
		repoClient := getRepoServiceClient(conn)
		respJson, err := callRepoMethod(repoClient, c)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, string(respJson))
	default:
		return errors.New("unrecognized service")
	}
	return nil
}
