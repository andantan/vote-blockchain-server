ifeq ($(OS),Windows_NT) 
	CLEAR_COMMAND = @cls
else 
	CLEAR_COMMAND = @clear
endif

CLIENT_DIR = impulse-client

test-debug:
	@go test -v ./...

test-race:
	@go test ./... --race

test:
	@go test ./...

install:
	@echo "--- Installing blockchain-node dependencies ---"
	@go mod tidy
	@echo "--- Installing blockchain-node dependencies completed ---"

build: install
	@echo "--- Building blockchain-node application ---"
	@go build -o ./bin/node
	@echo "--- Building blockchain-node application completed ---"

run-blockchain-node: build
	@$(CLEAR_COMMAND)
	./bin/node

run-burst-client:
	@echo "--- Running burst-client application ---"
	@make -C $(CLIENT_DIR) run ARGS="$(ARGS)"

