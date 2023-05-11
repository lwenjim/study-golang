module github.com/lwenjim/study-golang/code/demo10/server-healthcheck/server

go 1.19

require google.golang.org/grpc v1.55.0

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

require github.com/lwenjim/study-golang/code/demo10/server-healthcheck/service v0.0.0

replace github.com/lwenjim/study-golang/code/demo10/server-healthcheck/service => ../service
