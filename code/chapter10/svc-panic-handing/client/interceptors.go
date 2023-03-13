package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

func loggingUnaryInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(
		ctx,
		method,
		req,
		reply,
		cc,
		opts...,
	)
	end := time.Now()
	log.Printf(
		"Method:%s, Duration:%s, Error:%v",
		method,
		end.Sub(start),
		err,
	)
	return err
}

type wrappedClientStream struct {
	grpc.ClientStream
	messageSent int
	messageRcvd int
}

func (receiver wrappedClientStream) SendMsg(m interface{}) error {
	log.Printf("Send msg called: %T", m)
	err := receiver.ClientStream.SendMsg(m)
	receiver.messageSent += 1
	return err
}

func (receiver wrappedClientStream) RecvMsg(m interface{}) error {
	log.Printf("Receive msg called: %T", m)
	err := receiver.ClientStream.RecvMsg(m)
	receiver.messageRcvd += 1
	return err
}

func (receiver wrappedClientStream) CloseSend() error {
	log.Println("CloseSend msg called")
	v := receiver.Context().Value(streamDurationContextKey{})

	if m, ok := v.(streamDurationContextValue); ok {
		log.Printf("Duration:%v", time.Since(m.startTime))
	}
	err := receiver.ClientStream.CloseSend()
	log.Printf(
		"Messages Sent: %d, Messages Received:%d\n",
		receiver.messageSent,
		receiver.messageRcvd,
	)
	return err
}

type streamDurationContextKey struct {
}
type streamDurationContextValue struct {
	startTime time.Time
}

func loggingStreamingInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	c := streamDurationContextValue{
		startTime: time.Now(),
	}
	ctxWithTimer := context.WithValue(
		ctx,
		streamDurationContextKey{},
		c,
	)
	stream, err := streamer(
		ctxWithTimer,
		desc,
		cc, method,
		opts...,
	)
	clientStream := wrappedClientStream{
		ClientStream: stream,
		messageSent:  0,
		messageRcvd:  0,
	}
	return clientStream, err
}

func metadataUnaryInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,

) error {
	ctxWithMethodData := metadata.AppendToOutgoingContext(
		ctx,
		"Request-Id",
		"request-123",
	)
	return invoker(
		ctxWithMethodData,
		method,
		req,
		reply,
		cc,
		opts...,
	)
}
func metadataStreamInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (
	grpc.ClientStream, error) {
	ctxWithMetadata := metadata.AppendToOutgoingContext(
		ctx,
		"Request-Id",
		"request-123",
	)
	clientStream, err := streamer(
		ctxWithMetadata,
		desc,
		cc,
		method,
		opts...,
	)
	return clientStream, err
}
