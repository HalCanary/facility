PKGS=$(shell go list -f '{{.Dir}}' ./... | sed "s@^${PWD}/@@")
FILES=$(shell find . -name '*.go')
test: $(FILES)
	go get ./...
	go test ./...
fmt: $(FILES)
	find . -type f -name '*.go' -exec gofmt -w {} \;
update_deps:
	go get -u ./...
	go mod tidy
new_version:
	git tag $(shell V=$$(git describe --abbrev=0 --tags); echo $${V%.*}.$$(( $${V##*.}+1)))
push_up:
	git push --tags origin main:main
README.md: $(FILES) Makefile
	printf '# FACILITY\n\nUseful Go Packages\n\nCopyright (c) 2022 Hal Canary\n\n' > $@
	for pkg in $(PKGS); do \
		go doc -all ./$$pkg > $$pkg/README.txt; \
		printf '\n[%s](%s)\n' $$pkg ./$$pkg/ >> $@; \
	done
.PHONY: test fmt update_deps new_version push_up
