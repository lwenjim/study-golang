package main

import (
	"github.com/lwenjim/code/chapter10/exercise1/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"
)

var h *health.Server

func startTestGrpcServer() *bufconn.Listener {
	h = health.NewServer()
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	registerServices(s, h)
	updateServiceHealth(
		h,
		service.Users_ServiceDesc.ServiceName,
		grpc_health_v1.HealthCheckResponse_SERVING,
	)
	go func() {
		s.Serve(l)
	}()
	return l
}
