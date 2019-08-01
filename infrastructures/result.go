package infrastructures

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

type Result struct {
	IsSuccess bool
	Value     interface{}
	ErrorType int
	Error     error
}

const (
	ServerSideError = iota
	ClientSideError
	NotFound
	NoContent
)

func (result *Result) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	if result.IsSuccess {
		rw.WriteHeader(200)
		if err := producer.Produce(rw, result.Value); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	} else {
		rw.WriteHeader(400)
		if err := producer.Produce(rw, &struct{ Message string }{Message: result.Error.Error()}); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}

}
