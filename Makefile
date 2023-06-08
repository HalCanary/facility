test:
	go test ./...
fmt:
	find . -type f -name '*.go' -exec gofmt -w {} \;
.PHONY: test fmt
