package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"io/ioutil"
	stdlog "log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/mitchell/selfpass/credentials/middleware"
	"github.com/mitchell/selfpass/credentials/protobuf"
	"github.com/mitchell/selfpass/credentials/repositories"
	"github.com/mitchell/selfpass/credentials/service"
	"github.com/mitchell/selfpass/credentials/transport"
	"github.com/mitchell/selfpass/credentials/types"
)

var logger log.Logger

func main() {
	var (
		stop    = make(chan os.Signal, 1)
		dev     = flag.Bool("dev", false, "enables dev mode logging")
		port    = flag.String("port", "8080", "specify the port to listen on")
		crtFile = flag.String("cert", "/run/secrets/server_cert", "specify the cert file")
		keyFile = flag.String("key", "/run/secrets/server_key", "specify the private key file")
		caFile  = flag.String("ca", "/run/secrets/ca_cert", "specify the ca cert file")
		verbose = flag.Bool("v", false, "be more verbose")
		// tableName = flag.String(
		// 	"credential-table-name",
		// 	"selfpass-credential",
		// 	"specify the credential table name on AWS",
		// )
	)
	flag.Parse()

	signal.Notify(stop, syscall.SIGINT)
	signal.Notify(stop, syscall.SIGKILL)
	signal.Notify(stop, syscall.SIGTERM)

	logger = newLogger(os.Stdout, *dev)

	keypair, err := tls.LoadX509KeyPair(*crtFile, *keyFile)
	check(err)

	ca, err := ioutil.ReadFile(*caFile)
	check(err)

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{keypair},
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	// db := repositories.NewDynamoTable(*tableName)
	db, err := repositories.NewRedisConn("tcp", "redis:6379", 2)
	check(err)

	var svc types.Service
	svc = service.NewCredentials(db)
	if *verbose {
		svc = middleware.NewServiceLogger(logger, svc)
	}

	gsrv := transport.NewGRPCServer(svc, logger)
	srv := grpc.NewServer(grpc.Creds(creds))
	protobuf.RegisterCredentialServiceServer(srv, gsrv)

	addr := ":" + *port
	lis, err := net.Listen("tcp", addr)
	check(err)

	_ = logger.Log(
		"message", "serving",
		"address", addr,
		"dev", dev,
	)

	go func() { check(srv.Serve(lis)) }()

	<-stop
	_ = logger.Log("message", "gracefully stopping")
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
	l = log.WithPrefix(l, "caller", log.Caller(5), "timestamp", log.DefaultTimestamp)

	lfunc := log.LoggerFunc(func(in ...interface{}) error {
		if err := l.Log(in...); err != nil {
			stdlog.Println(err)
		}
		return nil
	})

	return lfunc
}

func check(err error) {
	if err != nil {
		_ = logger.Log("error", err)
		os.Exit(1)
	}
}
