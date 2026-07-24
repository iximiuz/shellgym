GO ?= go
BIN := bin/shellgym

.PHONY: build test vet run validate dist clean

build:
	CGO_ENABLED=0 $(GO) build -o $(BIN) ./cmd/shellgym

test:
	$(GO) test ./... -count=1 -timeout 300s

vet:
	$(GO) vet ./...

validate: build
	./$(BIN) validate --path paths/sample-linux-101

# Run the daemon against the reference path (playground use only).
run: build
	sudo ./$(BIN) serve --path paths/sample-linux-101 --addr :63636

# Assemble the self-provisioning playground bundle. The e2e playground's
# init task (e2e/playground.yaml) downloads it from
# https://labs.iximiuz.com/__static__/shellgym-dist.tar.gz - publish
# dist/shellgym-dist.tar.gz there after building.
dist: build
	rm -rf dist
	mkdir -p dist/shellgym/bin
	cp $(BIN) dist/shellgym/bin/shellgym
	cp -r paths dist/shellgym/paths
	install -m 0755 e2e/start.sh dist/shellgym/start.sh
	tar -C dist -czf dist/shellgym-dist.tar.gz shellgym

clean:
	rm -rf bin dist
