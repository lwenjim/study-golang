proto:
	@protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --proto_path=.  *.proto
	@protoc-go-inject-tag -input="./*.pb.go"
	@protofmt -w *.proto