.PHONY: all build clean format test gen-certs-go

build: clean gen-certs-go format
	go build -mod=vendor -o ./bin/server ./cmd/server
	rm ./cmd/server/certs.go

clean:
	rm -rf ./bin

docker: install
	docker-compose build

local:
	docker-compose up -d

up:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up

upd:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

down:
	docker-compose down

machine-create-google:
	docker-machine create --driver google \
	    --google-address selfpass \
	    --google-project selfpass-241808 \
	    --google-machine-type f1-micro \
	    --google-machine-image https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-9-stretch-v20190514 \
	    --google-username selfpass \
	    --google-zone us-west1-c \
	    selfpass01

machine-rm:
	docker-machine rm selfpass01

format:
	gofmt -w -s -l .

install:
	go mod tidy
	go mod vendor

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

gen-certs-go:
	./gen_certs_go.sh > ./cmd/server/certs.go

test:
	go test -cover ./...
