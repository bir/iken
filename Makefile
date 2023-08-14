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
