module github.com/mitchell/selfpass/services

go 1.12

require (
	github.com/etcd-io/bbolt v1.3.3
	github.com/go-kit/kit v0.9.0
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/mediocregopher/radix/v3 v3.3.0
	github.com/mitchell/selfpass/protobuf/go v0.0.0-00010101000000-000000000000
	github.com/spf13/pflag v1.0.3
	go.etcd.io/bbolt v1.3.3 // indirect
	google.golang.org/grpc v1.22.0
)

replace github.com/mitchell/selfpass/protobuf/go => ../protobuf/go
