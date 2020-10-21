.PHONY: build
build:
		go build -v cmd/dione/dione.go

test:
		go test -v -race -timeout 30s ./ ...

.DEFAULT_GOAL := build