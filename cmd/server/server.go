package main

import (
	"flag"
	"io"
	"net"
	"os"
	"os/signal"

	"github.com/go-kit/kit/log"
	"github.com/mitchell/selfpass/credentials/protobuf"
	"github.com/mitchell/selfpass/credentials/repositories"
	"github.com/mitchell/selfpass/credentials/service"
	"github.com/mitchell/selfpass/credentials/transport"
	"google.golang.org/grpc"
)

var logger log.Logger

func main() {
	var (
		stop      = make(chan os.Signal)
		dev       = flag.Bool("dev", false, "enables dev mode logging")
		port      = flag.String("port", "8080", "specify the port to listen on")
		tableName = flag.String(
			"credential-table-name",
			"selfpass-credential",
			"specify the credential table name on AWS",
		)
	)

	signal.Notify(stop, os.Interrupt)
	flag.Parse()

	logger = newLogger(os.Stdout, *dev)

	var (
		db   = repositories.NewDynamoTable(*tableName)
		svc  = service.NewCredentials(db)
		gsrv = transport.NewGRPCServer(svc, logger)
		srv  = grpc.NewServer()
	)
	protobuf.RegisterCredentialServiceServer(srv, gsrv)

	addr := "0.0.0.0:" + *port
	lis, err := net.Listen("tcp", addr)
	check(err)

	go func() {
		logger.Log(
			"message", "serving",
			"address", addr,
			"credentialTable", tableName,
			"dev", dev,
		)
		check(srv.Serve(lis))
	}()

	<-stop
	logger.Log("message", "gracefully stopping")
	srv.GracefulStop()
}

func newLogger(writer io.Writer, dev bool) log.Logger {
	var l log.Logger
	writer = log.NewSyncWriter(writer)

	if dev {
		l = log.NewLogfmtLogger(writer)
	} else {
		l = log.NewJSONLogger(writer)
	}
	l = log.WithPrefix(l, "caller", log.DefaultCaller, "timestamp", log.DefaultTimestamp)

	lfunc := log.LoggerFunc(func(in ...interface{}) error {
		if err := l.Log(in...); err != nil {
			panic(err.Error())
		}
		return nil
	})

	return lfunc
}

func check(err error) {
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}
}
