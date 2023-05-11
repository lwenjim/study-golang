package main

import (
	"context"
	"errors"
	"github.com/lwenjim/study-golang/code/demo10/client-resiliency/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"testing"
	"time"
)

func TestClientDxnInterceptor(t *testing.T) {
	req := service.UserGetRequest{}
	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "Users.GetUser",
	}
	testUnaryHandler := func(
		ctx context.Context,
		req interface{},
	) (interface{}, error) {
		time.Sleep(200 * time.Millisecond)
		return service.UserGetReply{}, nil
	}
	incomingContext, cancel := context.WithTimeout(
		context.Background(),
		100*time.Millisecond,
	)
	defer cancel()

	_, err := clientDisconnectUnaryInterceptor(
		incomingContext,
		&req,
		unaryInfo,
		testUnaryHandler,
	)
	expectedErr := status.Errorf(
		codes.Canceled,
		"Users.GetUser: Request canceled",
	)
	if !errors.Is(err, expectedErr) {
		t.Errorf(
			"Expected error:%v Got: %v\n",
			expectedErr,
			err,
		)
	}
}

type testStream struct {
	CancelFunc context.CancelFunc
	grpc.ServerStream
}

func (receiver testStream) Context() context.Context {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		100*time.Millisecond,
	)
	receiver.CancelFunc = cancel
	return ctx
}

func TestStreamingClientDxnInterceptor(t *testing.T) {
	streamInfo := &grpc.StreamServerInfo{
		FullMethod:     "Users.GetUser",
		IsClientStream: true,
		IsServerStream: true,
	}

	testStream := testStream{}
	testHandler := func(
		srv interface{},
		stream grpc.ServerStream,
	) (err error) {
		time.Sleep(200 * time.Millisecond)
		for {
			m := service.UserHelpRequest{}
			err := stream.RecvMsg(&m)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			r := service.UserHelpReply{}
			err = stream.SendMsg(&r)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := clientDisconnectStreamInterceptor(
		"test",
		testStream,
		streamInfo,
		testHandler,
	)
	expectedErr := status.Errorf(
		codes.Canceled,
		"Users.GetUser: Request canceled",
	)
	if !errors.Is(err, expectedErr) {
		t.Errorf(
			"Expected error: %v Got: %v\n",
			expectedErr,
			err,
		)
	}
}
