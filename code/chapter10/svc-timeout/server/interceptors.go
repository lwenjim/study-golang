package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type wrappedServerStream struct {
	grpc.ServerStream
	messageSent    int
	messageRcvd    int
	RecvMsgTimeout time.Duration
}

func (receiver wrappedServerStream) SendMsg(m interface{}) error {
	log.Printf("Send msg called:%T", m)
	err := receiver.ServerStream.SendMsg(m)
	receiver.messageSent += 1
	return err
}

func (receiver wrappedServerStream) RecvMsg(m interface{}) error {
	log.Printf("receive msg called:%T", m)
	receiver.messageRcvd += 1

	ch := make(chan error)
	t := time.NewTimer(receiver.RecvMsgTimeout)

	go func() {
		log.Printf(
			"Waiting to receive a message: %T",
			m,
		)
		ch <- receiver.ServerStream.RecvMsg(m)
	}()
	select {
	case <-t.C:
		return status.Error(
			codes.DeadlineExceeded,
			"DeadlineExceeded",
		)
	case err := <-ch:
		return err
	}
}

func timeoutUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (
	interface{},
	error,
) {
	var resp interface{}
	var err error
	ctxWithTimeout, cancel := context.WithTimeout(
		ctx,
		500*time.Millisecond,
	)
	defer cancel()
	ch := make(chan error)
	go func() {
		resp, err = handler(ctxWithTimeout, req)
		ch <- err
	}()
	select {
	case <-ctxWithTimeout.Done():
		cancel()
		err = status.Error(
			codes.DeadlineExceeded,
			fmt.Sprintf(
				"%s: DeadlineExceeded",
				info.FullMethod,
			),
		)
		return resp, err
	case <-ch:
	}
	return resp, err
}

func timeoutStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	serverStream := wrappedServerStream{
		RecvMsgTimeout: 500 * time.Millisecond,
		ServerStream:   stream,
	}
	err := handler(srv, serverStream)
	return err
}
