grpcurl -cacert ./certs/ca.pem \
        -cert ./certs/client.pem \
        -key ./certs/client-key.pem \
        -proto ./credentials/protobuf/service.proto \
        localhost:8080 \
        selfpass.credentials.CredentialService/Dump
