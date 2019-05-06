.PHONY: all build clean format test docker-build

build: clean format test
	go build -o ./bin/server ./cmd/server/server.go

docker-build:
	docker build -t selfpass .

clean:
	rm -rf ./bin
	go mod tidy

dev: docker-build
	docker run -i -t -p 8080:8080 selfpass -v -dev

local: docker-build
	docker run -i -t -p 8080:8080 selfpass

format:
	goimports -w -l .

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
