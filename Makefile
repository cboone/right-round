VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BINARY  := right-round
OUTDIR  := bin

LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build install clean lint vet test run tidy

build:
	mkdir -p $(OUTDIR)
	go build $(LDFLAGS) -o $(OUTDIR)/$(BINARY) ./cmd/right-round

install:
	go install $(LDFLAGS) ./cmd/right-round

clean:
	rm -rf $(OUTDIR) dist coverage.out

lint:
	golangci-lint run ./...

vet: ## Run go vet
	go vet ./...

test:
	go test ./... -race -coverprofile=coverage.out
	go tool cover -func=coverage.out

run: build
	./$(OUTDIR)/$(BINARY)

tidy:
	go mod tidy
