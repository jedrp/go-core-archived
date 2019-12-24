package apicore

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jedrp/go-core/pllog"
	flags "github.com/jessevdk/go-flags"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	schemeHTTP  = "http"
	schemeHTTPS = "https"
	schemeUnix  = "unix"
)

var defaultSchemes []string

func init() {
	defaultSchemes = []string{
		schemeHTTP,
	}
}

type Server interface {
	HTTPListener() (net.Listener, error)
	UnixListener() (net.Listener, error)
	TLSListener() (net.Listener, error)
	SetHandler(handler http.Handler)
	GetPort() int
	GetHandler() http.Handler
	GetTLSCertificate() string
	GetTLSCertificateKey() string
	GetEnabledListeners() []string
	Serve() (err error)
}

type CoreServer struct {
	Server
	DisableRest      bool   `long:"disable-rest" description:"Enable REST protocol" env:"DISABLE_REST"`
	DisableGrpc      bool   `long:"disable-grpc" description:"Enable Grpc protocol" env:"DISABLE_GRPC"`
	GrpcPort         int    `long:"grpc-port" description:"Enable Grpc protocol" env:"GRPC_PORT"`
	EnabledListener  string `long:"scheme" description:"Enable Grpc protocol" env:"SCHEME"`
	logger           pllog.PlLogger
	grpcServer       *grpc.Server
	enbaleTLSSetting bool
}
type ConfigGrpcFunc func(*grpc.Server)

func NewCoreServer(ctx context.Context,
	server Server,
	logger pllog.PlLogger,
	configGrpcServer ConfigGrpcFunc,
) *CoreServer {
	// set up REST server
	server.SetHandler(HandlePanicMiddleware(server.GetHandler(), logger))

	// set up gRPC server
	cert := server.GetTLSCertificate()
	certKey := server.GetTLSCertificateKey()
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
	if cert != "" || certKey != "" {
		creds, err := credentials.NewServerTLSFromFile(cert, certKey)
		if err != nil {
			logger.Fatal(err)
		}
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	grpcServer = grpc.NewServer(grpcOpts...)
	if configGrpcServer != nil {
		configGrpcServer(grpcServer)
	}

	coreServer := &CoreServer{
		Server:     server,
		logger:     logger,
		grpcServer: grpcServer,
	}

	if cert != "" || certKey != "" {
		coreServer.enbaleTLSSetting = true
	}

	parser := flags.NewParser(coreServer, flags.IgnoreUnknown)
	ParseConfig(parser)

	return coreServer
}

func (s *CoreServer) StartServing(ctx context.Context) error {

	s.logger.Infof("Sever start with REST_PORT: %v; GRPC_PORT: %v; SCHEME: %s; DISABLE_REST: %v; DISABLE_GRPC: %v;",
		s.GetPort(),
		s.GrpcPort,
		s.EnabledListener,
		s.DisableRest,
		s.DisableGrpc,
	)

	var l net.Listener
	var e error
	if s.EnabledListener == "http" {
		l, e = s.HTTPListener()
	} else {
		l, e = s.TLSListener()
	}

	if e != nil {
		log.Fatalf("failed to serve: %s", e.Error())
	}

	// serving both GRPC and REST
	if !s.DisableRest && !s.DisableGrpc {
		// same port for GRPC and REST, the we need special set up
		hostingGRPCAndRESTOnSamePort := s.GrpcPort > 0 && s.GetPort() == s.GrpcPort
		if hostingGRPCAndRESTOnSamePort {
			var server *http.Server
			if s.enbaleTLSSetting {
				cert, err := tls.LoadX509KeyPair(s.GetTLSCertificate(), s.GetTLSCertificateKey())
				if err != nil {
					panic(err)
				}
				tlsConfig := &tls.Config{
					Certificates: []tls.Certificate{cert},
				}
				l = tls.NewListener(l, tlsConfig)
				server = &http.Server{
					TLSConfig: tlsConfig,
					Handler:   http.HandlerFunc(s.getHandlerFunc()),
				}
				http2.ConfigureServer(server, nil)
			} else {
				server = &http.Server{
					Handler: h2c.NewHandler(s.getHandlerFunc(), &http2.Server{}),
				}
			}

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			go func() {
				for range c {
					s.logger.Info("shutting down server...")

					server.Close()

					<-ctx.Done()
				}
			}()

			s.logger.Infof("Sever starting serving both REST and gRPC at: %s\n", l.Addr())

			if err := server.Serve(l); !strings.Contains(err.Error(),
				"use of closed network connection") {
				s.logger.Fatal(err)
			}
			return nil

		} else {
			//spawn new thread to handle rest request
			go func() {
				s.logger.Infof("Sever starting serving REST at: %s\n", l.Addr())
				if err := s.Serve(); err != nil {
					log.Fatalln(err)
				}
			}()

			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.GrpcPort))
			if err != nil {
				s.logger.Panicf("failed to listen: %v", err)
			}
			s.logger.Infof("Sever starting serving gRPC at: %s\n", lis.Addr())
			if e = s.grpcServer.Serve(lis); e != nil {
				log.Fatalf("failed to serve: %s", e)
			}
		}
		return nil
	}

	// diable grpc
	if !s.DisableRest && s.DisableGrpc {
		fmt.Printf("Sever starting serving only REST at: %s\n", l.Addr())
		if err := s.Serve(); err != nil {
			log.Fatalln(err)
		}
	}

	//disable rest
	if s.DisableRest && !s.DisableGrpc {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.GrpcPort))
		if err != nil {
			s.logger.Panicf("failed to listen: %v", err)
		}
		fmt.Printf("Sever starting serving only gRPC at: %s\n", lis.Addr())
		if e = s.grpcServer.Serve(lis); e != nil {
			log.Fatalf("failed to serve: %s", e)
		}
	}

	return nil
}

func StartServers(ctx context.Context,
	server Server,
	logger pllog.PlLogger,
	configGrpcServer ConfigGrpcFunc,
) error {

	s := NewCoreServer(
		ctx,
		server,
		logger,
		configGrpcServer,
	)

	return s.StartServing(ctx)
}

func hasScheme(schemes []string, scheme string) bool {
	if len(schemes) == 0 {
		schemes = defaultSchemes
	}

	for _, v := range schemes {
		if v == scheme {
			return true
		}
	}
	return false
}

func ParseConfig(parser *flags.Parser) {
	if _, err := parser.Parse(); err != nil {
		log.Fatalln(err)
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}
}

func (s *CoreServer) getHandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			s.grpcServer.ServeHTTP(w, r)
		} else {
			s.GetHandler().ServeHTTP(w, r)
		}
	})
}
