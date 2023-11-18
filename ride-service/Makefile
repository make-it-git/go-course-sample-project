generate:
	mkdir -p internal/generated/proto
	protoc --go_out=./internal/generated/proto \
    	--go-grpc_out=./internal/generated/proto \
    	api/proto/service.proto

install-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	brew install protobuf
