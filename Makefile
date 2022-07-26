LOCAL_BIN := $(CURDIR)/bin

bin-deps:
	mkdir -p $(LOCAL_BIN)

	ls $(LOCAL_BIN)/goimports &> /dev/null || \
		GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@v0.1.9

	ls $(LOCAL_BIN)/golangci-lint &> /dev/null || \
		GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.47.1

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

run:
	go mod tidy && \
	go run cmd/bot/main.go
.PHONY: run