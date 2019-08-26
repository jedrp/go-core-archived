package plresult

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

//OKResult return ok reposne
func OKResult(value interface{}) *Result {
	return &Result{
		IsSuccess: true,
		Value:     value,
		Error:     nil,
	}
}

//ValidationErrorResult the first option param will be the code, the second one is the error message
func ValidationErrorResult(err error, opts ...string) *Result {
	errWrapper := newErrorResult(&ValidationError{}, err, opts)
	return &Result{
		IsSuccess: false,
		Value:     nil,
		Error:     errWrapper,
	}
}

//InternalErrorResult the first option param will be the code, the second one is the error message
func InternalErrorResult(err error, opts ...string) *Result {
	errWrapper := newErrorResult(&InternalServerError{}, err, opts)
	return &Result{
		IsSuccess: false,
		Value:     nil,
		Error:     errWrapper,
	}
}

//NotFoundErrorResult the first option param will be the code, the second one is the error message
func NotFoundErrorResult(err error, opts ...string) *Result {
	errWrapper := newErrorResult(&NotFoundError{}, err, opts)
	return &Result{
		IsSuccess: false,
		Value:     nil,
		Error:     errWrapper,
	}
}

func newErrorResult(errWrapper plError, err error, opts []string) plError {
	optsLength := len(opts)
	errWrapper.SetError(err)
	switch {
	case optsLength > 1:
		errWrapper.SetCode(opts[0])
		errWrapper.SetMessage(opts[1])
	case optsLength > 0:
		errWrapper.SetCode(opts[0])
		errWrapper.SetMessage(err.Error())
	default:
		errWrapper.SetCode("UNKNOWN_ERROR")
		errWrapper.SetMessage(err.Error())
	}
	return errWrapper
}

//WriteResponse append response to request
func (result *Result) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	if result.IsSuccess {
		rw.WriteHeader(200)
		if err := producer.Produce(rw, result.Value); err != nil {
			panic(err)
		}
	} else {
		err := result.Error
		responseMessage := &struct {
			Message   string
			ErrorCode string
		}{
			Message:   err.GetErrorMessage(),
			ErrorCode: err.GetCode(),
		}
		switch result.Error.(type) {
		case *ValidationError:
			rw.WriteHeader(400)
		case *NotFoundError:
			rw.WriteHeader(404)
		default:
			rw.WriteHeader(500)
		}
		if err := producer.Produce(rw, responseMessage); err != nil {
			panic(err)
		}

	}
}
