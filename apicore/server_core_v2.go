package apicore

import (
	"context"
	"net/http"
	"time"

	strfmt "github.com/go-openapi/strfmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jedrp/go-core/pllog"
	flags "github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type CoreServerV2 struct {
	GRPCPort          int            `long:"grpc-port" description:"Enable Grpc protocol" env:"GRPC_PORT"`
	RESTPort          int            `long:"rest-port" description:"Enable Grpc protocol" env:"REST_PORT"`
	ListenLimit       int            `long:"listen-limit" description:"limit the number of outstanding requests"`
	TLSCertificate    flags.Filename `long:"tls-certificate" description:"the certificate to use for secure connections" env:"TLS_CERTIFICATE"`
	TLSCertificateKey flags.Filename `long:"tls-key" description:"the private key to use for secure connections" env:"TLS_PRIVATE_KEY"`
	TLSCACertificate  flags.Filename `long:"tls-ca" description:"the certificate authority file to be used with mutual tls auth" env:"TLS_CA_CERTIFICATE"`
	KeepAlive         time.Duration  `long:"keep-alive" description:"sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)" default:"3m"`
	ReadTimeout       time.Duration  `long:"read-timeout" description:"maximum duration before timing out read of the request" default:"30s"`
	WriteTimeout      time.Duration  `long:"write-timeout" description:"maximum duration before timing out write of the response" default:"60s"`

	listenScheme     string
	restHandler      http.Handler
	logger           pllog.PlLogger
	grpcServer       *grpc.Server
	enbaleTLSSetting bool
}

func NewCoreServerV2(ctx context.Context,
	logger pllog.PlLogger,
	restHandler http.Handler,
	configGrpcServer ConfigGrpcFunc,
) *CoreServerV2 {

	coreServer := &CoreServerV2{
		logger:       logger,
		listenScheme: schemeHTTP,
	}

	parser := flags.NewParser(coreServer, flags.IgnoreUnknown)
	ParseConfig(parser)

	// set up REST server
	coreServer.restHandler = HandlePanicMiddleware(restHandler, logger)

	// set up gRPC server
	formats := strfmt.Default
	var grpcServer *grpc.Server

	grpcOpts := []grpc.ServerOption{grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		UnaryServerRequestContextInterceptor(),
		UnaryServerPanicInterceptor(logger),
		UnaryValidatorServerInterceptor(formats, logger),
	)), grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		StreamServerRequestInterceptor(),
		grpc_recovery.StreamServerInterceptor(
			grpc_recovery.WithRecoveryHandlerContext(getRecoveryHandlerFuncContextHandler(logger)),
		),
		StreamValidatorServerInterceptor(formats, logger),
	))}

	if coreServer.TLSCertificate != "" || coreServer.TLSCertificateKey != "" {
		coreServer.enbaleTLSSetting = true
		coreServer.listenScheme = schemeHTTPS

		creds, err := credentials.NewServerTLSFromFile(string(coreServer.TLSCertificate), string(coreServer.TLSCertificateKey))
		if err != nil {
			logger.Fatal(err)
		}
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	grpcServer = grpc.NewServer(grpcOpts...)
	if configGrpcServer != nil {
		configGrpcServer(grpcServer)
	}
	coreServer.grpcServer = grpcServer

	return coreServer
}

func (s *CoreServerV2) StartServing() {
	if s.GRPCPort < 1 && s.RESTPort < 1 {
		s.logger.Panicf("GRPC_PORT and REST_PORT are both not configured, stop!")
	}

	if s.GRPCPort < 1 && s.RESTPort > 1 {

	}

}
