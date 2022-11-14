test:
	go test -count=1 ./...

test-integration:
	go test -tags=integration -count=1 ./...

generate:
	go generate ./...

lint:
	golangci-lint run
