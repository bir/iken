lint:
	golangci-lint run --sort-results

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
	gofumpt -l -w .
	gci write . -s standard -s default -s "prefix(github.com/bir/iken)"

tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/daixiang0/gci@latest

.PHONY: lint test cover tidy update updateAll fmt tools

