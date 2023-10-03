.PHONY: lint test vendor clean

export GO111MODULE=on

default: lint test

lint:
	golangci-lint run

test:
	go test -count=1 -v -cover ./...

vendor:
	go mod vendor

clean:
	rm -rf ./vendor