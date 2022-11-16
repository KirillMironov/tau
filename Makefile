test:
	go test -count=1 ./...

test-integration:
	go test -tags=integration -count=1 ./...

generate:
	go generate ./...

proto:
	protoc --go_out=./api --go-grpc_out=require_unimplemented_servers=false:./api ./api/resources.proto
	go mod tidy

lint:
	golangci-lint run

tau-install:
	go install github.com/KirillMironov/tau/cmd/tau
