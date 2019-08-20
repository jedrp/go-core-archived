package platform

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

//Result Wrapper struct
type Result struct {
	IsSuccess bool
	Value     interface{}
	Error     plError
}

//WriteResponse append response to request
func (result *Result) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	if result.IsSuccess {
		rw.WriteHeader(200)
		if err := producer.Produce(rw, result.Value); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	} else {
		switch v := result.Error.(type) {
		case *ValidationError:
			rw.WriteHeader(400)
			if err := producer.Produce(rw, &struct{ Message string }{Message: v.GetErrorMessage()}); err != nil {
				panic(err) // let the recovery middleware deal with this
			}
		case *InternalServerError:
			rw.WriteHeader(500)
			if err := producer.Produce(rw, &struct{ Message string }{Message: result.Error.GetErrorMessage()}); err != nil {
				panic(err) // let the recovery middleware deal with this
			}
		default:
			rw.WriteHeader(400)
			if err := producer.Produce(rw, &struct{ Message string }{Message: result.Error.GetErrorMessage()}); err != nil {
				panic(err) // let the recovery middleware deal with this
			}
		}

	}
}
