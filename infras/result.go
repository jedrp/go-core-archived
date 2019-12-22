package infras

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Result struct {
	Value interface{}
	Error *status.Status
}

func OK(v interface{}) *Result {
	return &Result{
		v,
		nil,
	}
}

func Failf(code codes.Code, f string, o ...interface{}) *Result {
	return &Result{
		nil,
		status.Newf(code, f, o...),
	}
}
func Fail(code codes.Code, m string) *Result {
	return &Result{
		nil,
		status.New(code, m),
	}
}

// Implement Responder interface (Responder is an interface for types to implement, when they want to be considered for writing HTTP responses)
func (r *Result) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	if r.Error == nil {
		rw.WriteHeader(200)
		if err := producer.Produce(rw, r.Value); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	} else {
		code := r.Error.Code()
		switch code {
		case codes.InvalidArgument:
			rw.WriteHeader(400)
		case codes.NotFound:
			rw.WriteHeader(404)
		case codes.Aborted:
			rw.WriteHeader(412)
		default:
			rw.WriteHeader(500)
		}
		if err := producer.Produce(rw, &struct{ Message string }{Message: r.Error.Message()}); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
