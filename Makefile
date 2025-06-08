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


# 블록체인 노드 실행
# make blockchain-node-run

# Bursting 클라이언트 실행
# 생성할 투표에 관한 정보는 "impulse-client/data/vote_data.json" 참고
# makemake burst-client-run args="-max ${생성할 투표 개수}"

install:
	@echo "--- Installing blockchain-node dependencies ---"
	@go mod tidy
	@echo "--- Installing blockchain-node dependencies completed ---"

build: install
	@echo "--- Building blockchain-node application ---"
	@go build -o ./bin/node
	@echo "--- Building blockchain-node application completed ---"

blockchain-node-run: build
	@$(CLEAR_COMMAND)
	./bin/node

burst-client-run:
	@echo "--- Running burst-client application ---"
	@make -C $(CLIENT_DIR) run ARGS="$(ARGS)"

