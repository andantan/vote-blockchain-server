ifeq ($(OS),Windows_NT) 
	CLEAR_COMMAND = @cls
else 
	CLEAR_COMMAND = @clear
endif

install:
	@go mod tidy

build: install
	@go build -o ./bin/impulse-client

run: build
	@$(CLEAR_COMMAND)
	./bin/impulse-client $(args)
