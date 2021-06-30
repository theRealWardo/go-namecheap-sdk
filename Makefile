default: format check lint test

format:
	go fmt ./...

check:
	go vet ./...

test:
	go test -v ./...

# Make sure you have installed golangci-lint CLI
# https://golangci-lint.run/usage/install/#local-installation
lint:
	golangci-lint run
