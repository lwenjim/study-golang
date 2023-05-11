proto:
	@protoc --go_out=. --go-grpc_out=. --proto_path=.  *.proto
	@protoc-go-inject-tag -input="./*.pb.go"
	@protofmt -w *.proto