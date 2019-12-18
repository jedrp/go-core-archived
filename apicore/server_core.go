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
	"runtime/debug"
	"strings"

	"github.com/go-openapi/swag"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jedrp/go-core/pllog"
	flags "github.com/jessevdk/go-flags"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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
	ConfigureFlags()
	ConfigureAPI()
	HTTPListener() (net.Listener, error)
	UnixListener() (net.Listener, error)
	TLSListener() (net.Listener, error)
	SetHandler(handler http.Handler)
	GetHandler() http.Handler
	GetTLSCertificate() string
	GetTLSCertificateKey() string
	GetEnabledListeners() []string
	GetCommandLineOptionsGroups() []swag.CommandLineOptionsGroup
	Serve() (err error)
}

type CoreServer struct {
	Server
	logger          pllog.PlLogger
	grpcServer      *grpc.Server
	DisableRest     bool   `long:"disable-rest" description:"Enable REST protocol" env:"DISABLE_REST"`
	DisableGrpc     bool   `long:"disable-grpc" description:"Enable Grpc protocol" env:"DISABLE_GRPC"`
	EnabledListener string `long:"scheme" description:"Enable Grpc protocol" env:"SCHEME"`
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
	var grpcServer *grpc.Server
	grpcOpts := []grpc.ServerOption{grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		UnaryServerRequestContextInterceptor(),
		UnaryServerPanicInterceptor(logger),
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
	parser := flags.NewParser(coreServer, flags.IgnoreUnknown)
	ParseConfig(parser)

	return coreServer
}

func (s *CoreServer) StartServing(ctx context.Context) error {
	var schemes []string
	if s.EnabledListener != "" {
		schemes = []string{s.EnabledListener}
	} else {
		schemes = s.GetEnabledListeners()
	}
	var l net.Listener
	var e error
	if hasScheme(schemes, schemeHTTPS) {
		s.logger.Info("Server Https scheme enabled")
		l, e = s.TLSListener()
	} else if hasScheme(schemes, schemeHTTP) {
		s.logger.Info("Server Http scheme enabled")
		l, e = s.HTTPListener()
	} else if hasScheme(schemes, schemeUnix) {
		s.logger.Info("Server Unix scheme enabled")
		l, e = s.UnixListener()
	}
	if e != nil {
		log.Fatalf("failed to serve: %s", e)
	}

	if !s.DisableRest && !s.DisableGrpc {
		var server *http.Server
		if s.GetTLSCertificate() == "" || s.GetTLSCertificateKey() == "" {
			s.logger.Info(`----Sever starting serving both REST and gRPC without TLS setting (TLSCertificate and TLSCertificateKey) Using h2c to serve http2 insecure-----`)
			server = &http.Server{
				Handler: h2c.NewHandler(s.getHandlerFunc(), &http2.Server{}),
			}
		} else {
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

		fmt.Printf("Sever starting serving both REST and gRPC at: %s\n", l.Addr())
		if err := server.Serve(l); !strings.Contains(err.Error(),
			"use of closed network connection") {
			log.Fatal(err)
		}
		return nil
	}

	if !s.DisableRest && s.DisableGrpc {
		fmt.Printf("Sever starting serving only REST at: %s\n", l.Addr())
		if err := s.Serve(); err != nil {
			log.Fatalln(err)
		}
	}
	if s.DisableRest && !s.DisableGrpc {
		fmt.Printf("Sever starting serving only gRPC at: %s\n", l.Addr())
		if e = s.grpcServer.Serve(l); e != nil {
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

func UnaryServerPanicInterceptor(logger pllog.PlLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				if logger != nil {
					pllog.CreateLogEntryFromContext(ctx, logger).Error(r, string(debug.Stack()))
				}
			}
		}()

		return handler(ctx, req)
	}
}

func UnaryServerRequestContextInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {

		ctx = context.WithValue(ctx, pllog.RequestID, uuid.NewV4().String())

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			corID := md.Get(pllog.CorrelationIDHeaderKey)
			if len(corID) > 0 && corID[0] != "" {
				ctx = context.WithValue(ctx, pllog.CorrelationID, corID[0])
			}
			return handler(ctx, req)
		}

		err = fmt.Errorf("Unable to obtain metadata")
		return
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
