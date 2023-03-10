package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

type wrappedServerStream struct {
	grpc.ServerStream
	messageSent int
	messageRcvd int
}

func (receiver wrappedServerStream) SendMsg(m interface{}) error {
	log.Printf("Send msg called:%T", m)
	err := receiver.ServerStream.SendMsg(m)
	receiver.messageSent += 1
	return err
}

func (receiver wrappedServerStream) RecvMsg(m interface{}) error {
	log.Printf("receive msg called:%T", m)
	err := receiver.ServerStream.RecvMsg(m)
	receiver.messageRcvd += 1
	return err
}

func loggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (
	interface{},
	error,
) {
	start := time.Now()
	resp, err := handler(ctx, req)
	logMessage(
		context.Background(),
		info.FullMethod,
		time.Since(start),
		err,
	)
	return resp, err
}

func loggingStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	start := time.Now()
	err := handler(srv, stream)
	ctx := stream.Context()
	logMessage(
		ctx,
		info.FullMethod,
		time.Since(start),
		err,
	)
	return err
}

func metricUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	end := time.Now()
	log.Printf(
		"Method:%s, Duration:%s",
		info.FullMethod,
		end.Sub(start),
	)
	return resp, err
}

func metricStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	serverStream := wrappedServerStream{
		ServerStream: stream,
		messageSent:  0,
		messageRcvd:  0,
	}
	start := time.Now()
	err := handler(srv, serverStream)
	end := time.Now()
	log.Printf(
		"Method:%s, Duration:%s, Message Received:%d, Message Sent:%d",
		info.FullMethod,
		end.Sub(start),
		serverStream.messageRcvd,
		serverStream.messageSent,
	)
	return err
}
