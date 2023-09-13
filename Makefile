lint:
	golangci-lint run --sort-results

test:
	go test -race ./...

cover:
	go test	-coverprofile cp.out ./...
	go tool cover -html=cp.out

tidy:
	go mod tidy -compat=1.21

update: updateAll tidy

updateAll:
	go get -u ./...

fmt:
	gofumpt -l -w .
	gci write . -s standard -s default -s "prefix(github.com/bir/iken)"
