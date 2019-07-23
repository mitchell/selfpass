package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	stdlog "log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	protobuf "github.com/mitchell/selfpass/protobuf/go"
	"github.com/mitchell/selfpass/services/credentials/middleware"
	"github.com/mitchell/selfpass/services/credentials/repositories"
	"github.com/mitchell/selfpass/services/credentials/service"
	"github.com/mitchell/selfpass/services/credentials/transport"
	"github.com/mitchell/selfpass/services/credentials/types"
)

var logger log.Logger

func main() {
	var (
		stop     = make(chan os.Signal, 1)
		jsonLogs = flag.Bool("json-logs", false, "enables json logging")
		port     = flag.String("port", "8080", "specify the port to listen on")
		verbose  = flag.Bool("v", false, "be more verbose")
	)
	flag.Parse()

	signal.Notify(stop, syscall.SIGINT)
	signal.Notify(stop, syscall.SIGKILL)
	signal.Notify(stop, syscall.SIGTERM)

	logger = newLogger(os.Stdout, *jsonLogs)

	keypair, err := tls.X509KeyPair([]byte(cert), []byte(key))
	check(err)

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM([]byte(ca))

	creds := credentials.NewTLS(&tls.Config{
		Certificates:             []tls.Certificate{keypair},
		ClientCAs:                caPool,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
		},
	})

	db, err := repositories.OpenBoltDB("/home/selfpass/data/bolt.db", 0600, nil)
	check(err)

	var svc types.Service
	svc = service.NewCredentials(db)
	if *verbose {
		svc = middleware.NewServiceLogger(logger, svc)
	}

	gsrv := transport.NewGRPCServer(svc, logger)
	srv := grpc.NewServer(grpc.Creds(creds))
	protobuf.RegisterCredentialsServer(srv, gsrv)

	addr := ":" + *port
	lis, err := net.Listen("tcp", addr)
	check(err)

	_ = logger.Log(
		"message", "serving",
		"address", addr,
	)

	go func() { check(srv.Serve(lis)) }()

	<-stop
	_ = logger.Log("message", "gracefully stopping")
	srv.GracefulStop()
}

func newLogger(writer io.Writer, jsonLogs bool) log.Logger {
	writer = log.NewSyncWriter(writer)
	l := log.NewLogfmtLogger(writer)

	if jsonLogs {
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
