generate:
	mkdir -p internal/generated/proto
	protoc --go_out=./internal/generated/proto \
    	--go-grpc_out=./internal/generated/proto \
    	api/proto/service.proto

install-deps:
	go install github.com/jackc/tern/v2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	brew install protobuf

migrate:
	cd internal/db/migrations && tern migrate

up:
	DATABASE_URL="postgres://postgres:postgres@localhost:6433/postgres" \
	GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info \
	go run cmd/service/main.go
