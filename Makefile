GO ?= go
BIN := bin/shellgym

.PHONY: build test vet run validate clean

build:
	CGO_ENABLED=0 $(GO) build -o $(BIN) ./cmd/shellgym

test:
	$(GO) test ./... -count=1 -timeout 300s

vet:
	$(GO) vet ./...

validate: build
	./$(BIN) validate --content paths/sample-linux-101

# Run the daemon against the reference path (playground use only).
run: build
	sudo ./$(BIN) serve --content paths/sample-linux-101 --addr :63636

clean:
	rm -rf bin
