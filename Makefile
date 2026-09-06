.PHONY: test race cover vet fmt build tidy clean

test:
	go test -v -count=1 ./...

race:
	go test -race -v -count=1 ./... 

vet:
	go vet ./...

fmt:
	go fmt ./...

build:
	go build ./...

tidy:
	go mod tidy

clean:
	go clean ./...