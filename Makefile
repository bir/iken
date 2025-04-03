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
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/daixiang0/gci@latest

.PHONY: lint test cover tidy update updateAll fmt tools

