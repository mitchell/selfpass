.PHONY: all build clean format test

build: clean format test
	go build -o ./bin/server ./cmd/server/server.go

clean:
	rm -rf ./bin
	go mod tidy

dev: build
	./bin/server --dev

format:
	goimports -w -l .

gen-protoc:
	protoc --go_out=plugins=grpc:. \
		./credentials/protobuf/service.proto

test:
	go test -cover ./...
