test:
	go get ./...
	go test ./...
fmt:
	find . -type f -name '*.go' -exec gofmt -w {} \;
update_deps:
	go get -u ./...
	go mod tidy
new_version:
	git tag $(shell V=$$(git describe --abbrev=0 --tags); echo $${V%.*}.$$(( $${V##*.}+1)))
push_up:
	git push --tags origin main:main
.PHONY: test fmt update_deps new_version push_up
