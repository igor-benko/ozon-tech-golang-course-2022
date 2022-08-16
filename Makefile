LOCAL_BIN := $(CURDIR)/bin

bin-deps:
	mkdir -p $(LOCAL_BIN)

	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go install github.com/bufbuild/buf/cmd/buf@v1.6.0

lint: bin-deps
	$(CURDIR)/bin/golangci-lint run \
		--new-from-rev=origin/master \
		$(CURDIR)/...

tidy:
	go mod tidy

imports:
	find . \
		-name "*.go" \
		-not -path $(CURDIR)/.git \
		-not -path $(CURDIR)/vendor \
		-not -path $(GO_GENERATED_MYRIAD_FILES_DIR) \
		-not -path $(GO_GENERATED_SYSTEM_FILES_DIR) \
		-exec $(LOCAL_BIN)/goimports \
		-w {} \;

test:
	CGO_ENABLED=0 go test -v ./...

precommit: bin-deps tidy imports lint test

.PHONY: imports precommit test tidy bin-deps lint

run-bot:
	go mod tidy && \
	go build -o bot ./cmd/bot && \
	./bot
.PHONY: run-bot

run-server:
	go mod tidy && \
	go build -o server ./cmd/server && \
	./server
.PHONY: run-server

add-migration:
	goose -dir ./migrations create $(name) sql
.PHONY: add-migration