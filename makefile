.PHONY: test vet

check:
	go vet ./...
	go fmt ./...

test: check
	go test ./...

deps:
	dep ensure
