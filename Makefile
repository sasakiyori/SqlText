.PHONY: test

test:
	go test -race -covermode atomic -cover -v -coverprofile=coverage.out ./...
