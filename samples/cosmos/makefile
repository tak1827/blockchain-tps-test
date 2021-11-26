lint:
# 	${GOPATH}/bin/golint ./...
	go vet ./...

fmt:
	gofmt -w -l .

test:
	go test ./... -v -race

build:
	go build -o recorder -gcflags '-m'

time:
	gtime -f '%P %Uu %Ss %er %MkB %C' "$@" ./recorder

