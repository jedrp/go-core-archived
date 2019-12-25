package apicore

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-openapi/runtime/flagext"
	strfmt "github.com/go-openapi/strfmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jedrp/go-core/pllog"
	flags "github.com/jessevdk/go-flags"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/net/netutil"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type CoreServerV2 struct {
	Host              string           `long:"host" description:"the IP to listen on" default:"localhost" env:"HOST"`
	GRPCPort          int              `long:"grpc-port" description:"Enable Grpc protocol" env:"GRPC_PORT"`
	RESTPort          int              `long:"rest-port" description:"Enable Grpc protocol" env:"REST_PORT"`
	ListenLimit       int              `long:"listen-limit" description:"limit the number of outstanding requests"`
	TLSCertificate    flags.Filename   `long:"tls-certificate" description:"the certificate to use for secure connections" env:"TLS_CERTIFICATE"`
	TLSCertificateKey flags.Filename   `long:"tls-key" description:"the private key to use for secure connections" env:"TLS_PRIVATE_KEY"`
	TLSCACertificate  flags.Filename   `long:"tls-ca" description:"the certificate authority file to be used with mutual tls auth" env:"TLS_CA_CERTIFICATE"`
	KeepAlive         time.Duration    `long:"keep-alive" description:"sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)" default:"3m"`
	ReadTimeout       time.Duration    `long:"read-timeout" description:"maximum duration before timing out read of the request" default:"30s"`
	WriteTimeout      time.Duration    `long:"write-timeout" description:"maximum duration before timing out write of the response" default:"60s"`
	CleanupTimeout    time.Duration    `long:"cleanup-timeout" description:"grace period for which to wait before killing idle connections" default:"10s"`
	GracefulTimeout   time.Duration    `long:"graceful-timeout" description:"grace period for which to wait before shutting down the server" default:"15s"`
	MaxHeaderSize     flagext.ByteSize `long:"max-header-size" description:"controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of the request body." default:"1MiB"`

	listenScheme     string
	restHandler      http.Handler
	logger           pllog.PlLogger
	grpcServer       *grpc.Server
	restServer       *http.Server
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

func (s *CoreServerV2) StartServing(ctx context.Context) error {
	if s.GRPCPort < 1 && s.RESTPort < 1 {
		s.logger.Panicf("GRPC_PORT and REST_PORT are both not configured, stop!")
	}

	if s.GRPCPort > 1 && s.RESTPort < 1 {
		// enable only grp
		return s.serveGRPCAPI()
	} else if s.GRPCPort < 1 && s.RESTPort > 1 {
		// enable only rest
		return s.serveRESTAPI()

	} else {

		if s.GRPCPort != s.RESTPort {
			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
			defer signal.Stop(interrupt)

			g, ctx := errgroup.WithContext(ctx)
			// run grpc and rest on difference ports
			g.Go(func() error {
				return s.serveGRPCAPI()
			})
			g.Go(func() error {
				return s.serveRESTAPI()
			})
			select {
			case <-interrupt:
				break
			case <-ctx.Done():
				break
			}
			if s.grpcServer != nil {
				s.grpcServer.GracefulStop()
			}
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()
			if s.restServer != nil {
				_ = s.restServer.Shutdown(shutdownCtx)
			}

		} else {
			// run grpc and rest in shared port mode
			s.serveGRPCAndRESTInSharedPortMode()
		}

	}
	return nil
}

func (s *CoreServerV2) getHandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			s.grpcServer.ServeHTTP(w, r)
		} else {
			s.restHandler.ServeHTTP(w, r)
		}
	})
}
func (s *CoreServerV2) serveGRPCAPI() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.GRPCPort))
	if err != nil {
		s.logger.Panicf("failed to listen: %v", err)
	}
	s.logger.Infof("Sever starting serving gRPC at: %s\n", lis.Addr())
	if err = s.grpcServer.Serve(lis); err != nil {
		s.logger.Fatalf("failed to serve: %s", err)
		return err
	}
	return nil
}

func (s *CoreServerV2) serveGRPCAndRESTInSharedPortMode() error {
	var server *http.Server
	l, err := net.Listen("tcp", net.JoinHostPort(s.Host, strconv.Itoa(s.RESTPort)))
	if err != nil {
		return err
	}
	if s.enbaleTLSSetting {
		cert, err := tls.LoadX509KeyPair(string(s.TLSCertificate), string(s.TLSCertificateKey))
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
		}
	}()
	s.logger.Infof("Sever starting serving both REST and gRPC at: %s\n", l.Addr())

	s.restServer = server

	if err := server.Serve(l); !strings.Contains(err.Error(),
		"use of closed network connection") {
		s.logger.Fatal(err)
		return err
	}
	return nil
}
func (s *CoreServerV2) serveRESTAPI() error {
	switch s.listenScheme {
	case "http":
		listener, err := net.Listen("tcp", net.JoinHostPort(s.Host, strconv.Itoa(s.RESTPort)))
		if err != nil {
			return err
		}
		httpServer := new(http.Server)
		httpServer.MaxHeaderBytes = int(s.MaxHeaderSize)
		httpServer.ReadTimeout = s.ReadTimeout
		httpServer.WriteTimeout = s.WriteTimeout
		httpServer.SetKeepAlivesEnabled(int64(s.KeepAlive) > 0)
		if s.ListenLimit > 0 {
			listener = netutil.LimitListener(listener, s.ListenLimit)
		}

		if int64(s.CleanupTimeout) > 0 {
			httpServer.IdleTimeout = s.CleanupTimeout
		}
		httpServer.Handler = s.restHandler
		s.restServer = httpServer
		s.logger.Infof("Sever starting serving REST(http)  at: %s", listener.Addr())
		if err := httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("%v", err)
			return err
		}

	case "https":
		tlsListener, err := net.Listen("tcp", net.JoinHostPort(s.Host, strconv.Itoa(s.RESTPort)))

		if err != nil {
			return err
		}
		httpsServer := new(http.Server)
		httpsServer.MaxHeaderBytes = int(s.MaxHeaderSize)
		httpsServer.ReadTimeout = s.ReadTimeout
		httpsServer.WriteTimeout = s.WriteTimeout
		httpsServer.SetKeepAlivesEnabled(int64(s.KeepAlive) > 0)
		if s.ListenLimit > 0 {
			tlsListener = netutil.LimitListener(tlsListener, s.ListenLimit)
		}
		if int64(s.CleanupTimeout) > 0 {
			httpsServer.IdleTimeout = s.CleanupTimeout
		}
		httpsServer.Handler = s.restHandler

		// Inspired by https://blog.bracebin.com/achieving-perfect-ssl-labs-score-with-go
		httpsServer.TLSConfig = &tls.Config{
			// Causes servers to use Go's default ciphersuite preferences,
			// which are tuned to avoid attacks. Does nothing on clients.
			PreferServerCipherSuites: true,
			// Only use curves which have assembly implementations
			// https://github.com/golang/go/tree/master/src/crypto/elliptic
			CurvePreferences: []tls.CurveID{tls.CurveP256},
			// Use modern tls mode https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
			NextProtos: []string{"h2", "http/1.1"},
			// https://www.owasp.org/index.php/Transport_Layer_Protection_Cheat_Sheet#Rule_-_Only_Support_Strong_Protocols
			MinVersion: tls.VersionTLS12,
			// These ciphersuites support Forward Secrecy: https://en.wikipedia.org/wiki/Forward_secrecy
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			},
		}
		if s.TLSCertificate == "" && s.TLSCertificateKey == "" {
			return fmt.Errorf("scheme %s enable required TLSCertificate and TLSCertificateKey cause server to stop", s.listenScheme)
		}
		// build standard config from server options
		httpsServer.TLSConfig.Certificates = make([]tls.Certificate, 1)
		httpsServer.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(string(s.TLSCertificate), string(s.TLSCertificateKey))
		if err != nil {
			return err
		}

		if s.TLSCACertificate != "" {
			// include specified CA certificate
			caCert, caCertErr := ioutil.ReadFile(string(s.TLSCACertificate))
			if caCertErr != nil {
				return caCertErr
			}
			caCertPool := x509.NewCertPool()
			ok := caCertPool.AppendCertsFromPEM(caCert)
			if !ok {
				return fmt.Errorf("cannot parse CA certificate")
			}
			httpsServer.TLSConfig.ClientCAs = caCertPool
			httpsServer.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
		if len(httpsServer.TLSConfig.Certificates) == 0 && httpsServer.TLSConfig.GetCertificate == nil {
			// after standard and custom config are passed, this ends up with no certificate
			if s.TLSCertificate == "" {
				if s.TLSCertificateKey == "" {
					s.logger.Fatalf("the required flags `--tls-certificate` and `--tls-key` were not specified")
				}
				s.logger.Fatalf("the required flag `--tls-certificate` was not specified")
			}
			if s.TLSCertificateKey == "" {
				s.logger.Fatalf("the required flag `--tls-key` was not specified")
			}
			// this happens with a wrong custom TLS configurator
			s.logger.Fatalf("no certificate was configured for TLS")
		}

		// must have at least one certificate or panics
		httpsServer.TLSConfig.BuildNameToCertificate()
		s.restServer = httpsServer
		s.logger.Infof("Sever starting serving REST(https)  at: %s", tlsListener.Addr())
		if err := httpsServer.Serve(tlsListener); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("%v", err)
			return err
		}
	default:
		return fmt.Errorf("scheme %s not supported", s.listenScheme)
	}
	return nil
}
