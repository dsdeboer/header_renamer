.PHONY: lint test_go test_yaegi vendor clean

export GO111MODULE=on

default: lint test

test: test_yaegi test_go

lint:
	golangci-lint run

test_go:
	go test -v -cover ./...

test_yaegi:
	yaegi test -v .

vendor:
	go mod vendor

clean:
	rm -rf ./vendor

