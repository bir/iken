lint:
	golangci-lint run

test:
	go test -race ./...

cover:
	go test	-coverprofile cp.out ./...
	go tool cover -html=cp.out

tidy:
	go mod tidy

update: updateAll tidy

updateAll:
	go get -u ./...

fmt:
	golangci-lint fmt

tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

.PHONY: lint test cover tidy update updateAll fmt tools

