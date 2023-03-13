package main

import (
	"context"
	"errors"
	"github.com/lwenjim/code/chapter10/svc-timeout/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"testing"
	"time"
)

func TestUnaryTimeoutInterceptor(t *testing.T) {
	req := service.UserGetRequest{}
	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "Users.GetUser",
	}
	testUnaryHandler := func(
		ctx context.Context,
		req interface{},
	) (
		interface{},
		error,
	) {
		time.Sleep(300 * time.Millisecond)
		return service.UserGetReply{}, nil
	}
	_, err := timeoutUnaryInterceptor(
		context.Background(),
		&req,
		unaryInfo,
		testUnaryHandler,
	)
	if err == nil {
		t.Fatal(err)
	}
	expectedErr := status.Errorf(
		codes.DeadlineExceeded,
		"Users.GetUser: DeadlineExceeded",
	)
	if !errors.Is(err, expectedErr) {
		t.Errorf(
			"Expected error: %v Got: %v\n",
			expectedErr,
			err,
		)
	}
}

type testStream struct {
	grpc.ServerStream
}

func (receiver testStream) SendMsg(m interface{}) error {
	log.Println("Test Stream - Sending")
	return nil
}

func (receiver testStream) RecvMsg(m interface{}) error {
	log.Println("Test Steram - RecvMsg - Going to sleep")
	time.Sleep(700 * time.Millisecond)
	return nil
}

func TestStreamingTimeOutInterceptor(t *testing.T) {
	streamInfo := &grpc.StreamServerInfo{
		FullMethod:     "Users.GetUser",
		IsClientStream: true,
		IsServerStream: true,
	}
	testStream := testStream{}

	testHandler := func(
		srv interface{},
		stream grpc.ServerStream,
	) error {
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
	err := timeoutStreamInterceptor(
		"test",
		testStream,
		streamInfo,
		testHandler,
	)
	expectedErr := status.Errorf(
		codes.DeadlineExceeded,
		"DeadlineExceeded",
	)
	if !errors.Is(err, expectedErr) {
		t.Errorf(
			"Expected error: %v Got: %v\n",
			expectedErr,
			err,
		)
	}
}
