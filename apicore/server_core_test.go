package apicore

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/jedrp/go-core/pllog"
	"google.golang.org/grpc"
)

type emptyServiceServer interface{}

type testServer struct{}

type testRestServer struct{}

func TestGetServiceInfo(t *testing.T) {
	testSd := grpc.ServiceDesc{
		ServiceName: "grpc.testing.EmptyService",
		HandlerType: (*emptyServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "EmptyCall",
				Handler:    nil,
			},
		},
		Streams: []grpc.StreamDesc{
			{
				StreamName:    "EmptyStream",
				Handler:       nil,
				ServerStreams: false,
				ClientStreams: true,
			},
		},
		Metadata: []int{0, 2, 1, 3},
	}

	restServer := &testRestServer{}
	ctx := context.Background()
	logger := &pllog.DefaultLogger{}

	s := NewCoreServer(ctx,
		restServer,
		logger,
		func(grpcServer *grpc.Server) {
			grpcServer.RegisterService(&testSd, &testServer{})
		},
	)
	if s == nil {
		t.Error("fail to start server")
	}
}

func (s *testRestServer) HTTPListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", net.JoinHostPort("localhost", strconv.Itoa(8080)))
	if err != nil {
		return nil, err
	}
	return listener, nil
}
func (s *testRestServer) UnixListener() (net.Listener, error) {
	return nil, nil
}
func (s *testRestServer) TLSListener() (net.Listener, error) {
	return nil, nil
}
func (s *testRestServer) SetHandler(handler http.Handler) {
	return
}
func (s *testRestServer) GetPort() int {
	return 0
}
func (s *testRestServer) GetHost() string {
	return "localhost"
}
func (s *testRestServer) GetHandler() http.Handler {
	return nil
}
func (s *testRestServer) GetTLSCertificate() string {
	return ""
}
func (s *testRestServer) GetTLSCertificateKey() string {
	return ""
}
func (s *testRestServer) GetEnabledListeners() []string {
	return []string{""}
}
func (s *testRestServer) GetCommandLineOptionsGroups() []swag.CommandLineOptionsGroup {
	return nil
}
func (s *testRestServer) Serve() (err error) {
	return nil
}
