.PHONY: all build clean format test docker-build

build: clean format test
	go build --o ./bin/server ./cmd/server/server.go

clean:
	rm -rf ./bin
	go mod tidy

docker:
	docker-compose build

start:
	docker-compose up

format:
	gofmt -w -s -l .

install-spc:
	go install ./cmd/spc

gen-protoc:
	protoc --go_out=plugins=grpc:. \
		./credentials/protobuf/service.proto

gen-csr-json:
	mkdir certs
	cd certs && cfssl print-defaults csr > csr.json

gen-ca:
	cd certs && cfssl genkey -initca csr.json | cfssljson -bare ca

gen-server-cert:
	cd certs && cfssl gencert -ca ca.pem -ca-key ca-key.pem csr.json | cfssljson -bare server

gen-client-cert:
	cd certs && cfssl gencert -ca ca.pem -ca-key ca-key.pem csr.json | cfssljson -bare client

test:
	go test -cover ./...
