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
		verbose  = flag.Bool("v", false, "be more verbose")
		jsonLogs = flag.Bool("json-logs", false, "enables json logging")
		port     = flag.String("port", "8080", "specify the port to listen on")
		caFile   = flag.String("ca", "/run/secrets/ca", "specify an alternate ca file")
		certFile = flag.String("cert", "/run/secrets/cert", "specify an alternate cert file")
		keyFile  = flag.String("key", "/run/secrets/key", "specify an alternate key file")
		boltFile = flag.String("bolt-file", "/home/selfpass/data/bolt.db", "specify an alternate bolt db file")
	)
	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	signal.Notify(stop, syscall.SIGTERM)

	logger = newLogger(os.Stdout, *jsonLogs)

	ca, err := ioutil.ReadFile(*caFile)
	check(err)
	cert, err := ioutil.ReadFile(*certFile)
	check(err)
	key, err := ioutil.ReadFile(*keyFile)
	check(err)

	keypair, err := tls.X509KeyPair(cert, key)
	check(err)

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates:             []tls.Certificate{keypair},
		ClientCAs:                caPool,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		MinVersion:               tls.VersionTLS13,
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
		},
	})

	srv := grpc.NewServer(grpc.Creds(creds))

	db, err := repositories.OpenBoltDB(*boltFile, 0600, nil)
	check(err)

	var svc types.Service
	svc = service.NewCredentials(db)
	if *verbose {
		svc = middleware.NewServiceLogger(logger, svc)
	}

	gsrv := transport.NewGRPCServer(svc, logger)
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
