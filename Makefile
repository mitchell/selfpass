.PHONY: all build clean format test docker-build

build: clean format
	env CGO_ENABLED=0 go build -o ./bin/server ./cmd/server
	rm ./cmd/server/certs.go

clean:
	rm -rf ./bin
	go mod tidy

docker:
	docker-compose build

up:
	docker-compose up

upd:
	docker-compose up -d

down:
	docker-compose down

machine-create-google:
	docker-machine create --driver google \
	    --google-address m-selfpass \
	    --google-project selfpass-241808 \
	    --google-machine-type n1-standard-1 \
	    --google-machine-image https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-9-stretch-v20190514 \
	    selfpass01

machine-rm:
	docker-machine rm selfpass01

machine-ssh:
	docker-machine ssh selfpass01

machine-env:
	docker-machine env selfpass01

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

gen-certs-go:
	./gen_certs_go.sh > ./cmd/server/certs.go

test:
	go test -cover ./...
