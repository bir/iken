lint:
	golangci-lint run --sort-results

test:
	go test -race ./...

cover:
	go test	-coverprofile cp.out ./...
	go tool cover -html=cp.out
