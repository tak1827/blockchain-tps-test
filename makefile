lint:
# 	${GOPATH}/bin/golint ./...
	go vet ./...

fmt:
	gofmt -w -l .

test:
	go test ./... -v -race
