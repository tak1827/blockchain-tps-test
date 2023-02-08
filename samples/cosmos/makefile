lint:
# 	${GOPATH}/bin/golint ./...
	go vet ./...

fmt:
	gofmt -w -l .

test:
	go test ./... -v -race

run:
	go run main

build:
	go build -o recorder -gcflags '-m'

time:
	gtime -f '%P %Uu %Ss %er %MkB %C' "$@" ./recorder

