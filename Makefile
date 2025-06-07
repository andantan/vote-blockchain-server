ifeq ($(OS),Windows_NT) 
	CLEAR_COMMAND = @cls
else 
	CLEAR_COMMAND = @clear
endif

gen-protobuf:
	@protoc \
	--proto_path=protobuf -I "proto/transaction.proto" \
	--go_out=protobuf \
	--go_opt=paths=source_relative \
    --go-grpc_out=./protobuf\
	--go-grpc_opt=paths=source_relative \

build:
	@go build -o ./bin/node

run: build
	@$(CLEAR_COMMAND)
	./bin/node

test-debug:
	@go test -v ./...

test-race:
	@go test ./... --race

test:
	@go test ./...