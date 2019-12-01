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

	"github.com/go-openapi/loads"
	"github.com/go-openapi/swag"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jedrp/go-core/pllog"
	flags "github.com/jessevdk/go-flags"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
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

type RestServer interface {
	ConfigureFlags()
	ConfigureAPI()
	HTTPListener() (net.Listener, error)
	UnixListener() (net.Listener, error)
	TLSListener() (net.Listener, error)
	GetHandler() http.Handler
	GetTLSCertificate() string
	GetTLSCertificateKey() string
	GetEnabledListeners() []string
	GetCommandLineOptionsGroups() []swag.CommandLineOptionsGroup
	Serve() (err error)
}

type CoreServer struct {
	restServer  RestServer
	grpcServer  *grpc.Server
	DisableRest bool `long:"disable-rest" description:"Enable REST protocol" env:"DISABLE_REST"`
	DisableGrpc bool `long:"disable-grpc" description:"Enable Grpc protocol" env:"DISABLE_GRPC"`
}
type ConfigGrpcFunc func(*grpc.Server)

func StartServers(ctx context.Context,
	restServer RestServer,
	logger pllog.PlLogger,
	configGrpcServer ConfigGrpcFunc,
	swaggerSpec *loads.Document,
) *CoreServer {

	// set up grpc server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		UnaryServerRequestContextInterceptor(),
		UnaryServerPanicInterceptor(logger),
	)))

	if configGrpcServer != nil {
		configGrpcServer(grpcServer)
	}

	// set up REST server
	parser := flags.NewParser(restServer, flags.IgnoreUnknown)
	parser.ShortDescription = swaggerSpec.Spec().Info.Title
	parser.LongDescription = swaggerSpec.Spec().Info.Description
	restServer.ConfigureFlags()
	for _, optsGroup := range restServer.GetCommandLineOptionsGroups() {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}
	ParseConfig(parser)
	restServer.ConfigureAPI()

	coreServer := &CoreServer{
		restServer: restServer,
		grpcServer: grpcServer,
	}

	parser = flags.NewParser(coreServer, flags.IgnoreUnknown)
	ParseConfig(parser)

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	schemes := restServer.GetEnabledListeners()
	fmt.Println("schemes", len(schemes))
	var l net.Listener
	var e error
	if hasScheme(schemes, schemeHTTPS) {
		l, e = restServer.TLSListener()
	} else if hasScheme(schemes, schemeHTTP) {
		l, e = restServer.HTTPListener()
	} else if hasScheme(schemes, schemeUnix) {
		l, e = restServer.UnixListener()
	}

	if !coreServer.DisableRest && !coreServer.DisableGrpc {
		l = tls.NewListener(l, tlsConfig)
		server := &http.Server{
			TLSConfig: tlsConfig,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
					grpcServer.ServeHTTP(w, r)
				} else {
					restServer.GetHandler().ServeHTTP(w, r)
				}
			}),
		}
		http2.ConfigureServer(server, nil)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for range c {
				// sig is a ^C, handle it
				logger.Info("shutting down gRPC server...")

				server.Close()

				<-ctx.Done()
			}
		}()
		// l, _ := restServer.HTTPListener()
		fmt.Printf("Sever starting serving both REST and gRPC at: %s\n", l.Addr())
		if err := server.Serve(l); !strings.Contains(err.Error(),
			"use of closed network connection") {
			log.Fatal(err)
		}
		return coreServer
	}

	if e != nil {
		log.Fatalf("failed to serve: %s", e)
	}
	if !coreServer.DisableRest && coreServer.DisableGrpc {
		fmt.Printf("Sever starting serving only REST at: %s\n", l.Addr())
		if err := restServer.Serve(); err != nil {
			log.Fatalln(err)
		}
	}
	if coreServer.DisableRest && !coreServer.DisableGrpc {
		fmt.Printf("Sever starting serving only gRPC at: %s\n", l.Addr())
		if e = grpcServer.Serve(l); e != nil {
			log.Fatalf("failed to serve: %s", e)
		}
		return nil
	}
	return nil
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
