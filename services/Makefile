.PHONY: all build clean format test

docker: clean format install
	docker build -t selfpass:latest .

build:
	go build -mod=vendor -o ./bin/server ./cmd

clean:
	rm -rf ./bin ./vendor

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
	    ${DOCKER_MACHINE_NAME}
	$(MAKE) machine-install-stackdriver-agent
	$(MAKE) machine-add-grpc-server-tag

machine-rm:
	docker-machine rm ${DOCKER_MACHINE_NAME}

machine-ssh:
	docker-machine ssh ${DOCKER_MACHINE_NAME}

machine-put-data:
	docker-machine ssh ${DOCKER_MACHINE_NAME} "if [[ ! -e data ]]; then mkdir data && chmod 777 data; fi"
	docker-machine scp ./data/bolt.db ${DOCKER_MACHINE_NAME}:data/bolt.db
	docker-machine ssh ${DOCKER_MACHINE_NAME} "chmod 666 data/bolt.db"

machine-get-data:
	docker-machine scp ${DOCKER_MACHINE_NAME}:data/bolt.db ./data/

machine-put-certs:
	docker-machine ssh ${DOCKER_MACHINE_NAME} "if [[ ! -e certs ]]; then mkdir certs; fi && chmod -R 755 certs"
	docker-machine scp ./certs/ca.pem ${DOCKER_MACHINE_NAME}:certs/ca.pem
	docker-machine scp ./certs/server.pem ${DOCKER_MACHINE_NAME}:certs/server.pem
	docker-machine scp ./certs/server-key.pem ${DOCKER_MACHINE_NAME}:certs/server-key.pem
	docker-machine ssh ${DOCKER_MACHINE_NAME} "chmod 444 certs/*"

machine-add-grpc-server-tag:
	gcloud compute instances add-tags ${DOCKER_MACHINE_NAME} \
		--zone us-west1-c \
		--tags grpc-server

machine-install-stackdriver-agent:
	docker-machine ssh ${DOCKER_MACHINE_NAME} "curl -sSO https://dl.google.com/cloudagents/install-monitoring-agent.sh && sudo bash install-monitoring-agent.sh"

format:
	gofmt -w -s -l .

install:
	go mod tidy
	go mod vendor

install-spc:
	go install ./cmd/spc

gen-protoc:
	protoc --go_out=plugins=grpc:. \
		--dart_out=grpc:. \
		./credentials/protobuf/service.proto

gen-csr-json:
	mkdir certs
	cd certs && cfssl print-defaults csr > csr.json

gen-ca:
	cd certs && cfssl genkey -initca csr.json | cfssljson -bare ca

gen-server-cert:
	cd certs && cfssl gencert -ca ca.pem -ca-key ca-key.pem -profile www csr.json | cfssljson -bare server

gen-client-cert:
	cd certs && cfssl gencert -ca ca.pem -ca-key ca-key.pem -profile client csr.json | cfssljson -bare client

test:
	go test -cover ./...
