build:
	@go build -o ./bin/node

run: build
	./bin/node


test-debug:
	@go test -v ./...

test:
	@go test ./...