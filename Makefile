.PHONY: all build clean format test gen-certs-go

docker: clean format install
	docker-compose build

build: gen-certs-go
	go build -mod=vendor -o ./bin/server ./cmd/server
	rm ./cmd/server/certs.go

clean:
	rm -rf ./bin ./vendor ./cmd/server/certs.go

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
	    --google-machine-type g1-small \
	    --google-machine-image https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-9-stretch-v20190514 \
	    --google-username selfpass \
	    --google-zone us-west1-c \
	    selfpass01
	$(MAKE) machine-put-redis.conf
	$(MAKE) machine-put-data
	$(MAKE) machine-install-stackdriver-agent
	$(MAKE) machine-add-grpc-server-tag

machine-rm:
	docker-machine rm selfpass01

machine-ssh:
	docker-machine ssh selfpass01

machine-put-redis.conf:
	docker-machine scp ./redis.conf selfpass01:redis.conf

machine-put-data:
	docker-machine scp -r ./data selfpass01:

machine-get-data:
	docker-machine scp -r selfpass01:data ./

machine-add-grpc-server-tag:
	gcloud compute instances add-tags selfpass01 \
		--zone us-west1-c \
		--tags grpc-server

machine-install-stackdriver-agent:
	docker-machine ssh selfpass01 "curl -sSO https://dl.google.com/cloudagents/install-monitoring-agent.sh && sudo bash install-monitoring-agent.sh"

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
