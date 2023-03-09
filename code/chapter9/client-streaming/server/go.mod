module github.com/lwenjim/code/chapter9/server-streaming/server

go 1.16

require (
	github.com/lwenjim/code/chapter9/server-streaming/service v0.0.0
	google.golang.org/grpc v1.53.0
)

replace github.com/lwenjim/code/chapter9/server-streaming/service => ../service
