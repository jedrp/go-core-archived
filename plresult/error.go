package plresult

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error interface {
	GetCode() string
	GetOriginError() error
	GetErrorMessage() string
	SetCode(code string)
	SetError(err error)
	SetMessage(msg string)
}

//Error wrapper
type errorObj struct {
	OriginError  error  //the origin error causing problem
	ErrorCode    string //System specific error code
	ErrorMessage string //System specfic error message
}

//ValidationError ..
type ValidationError errorObj

func (err *ValidationError) GetCode() string {
	return err.ErrorCode
}
func (err *ValidationError) GetOriginError() error {
	return err.OriginError
}
func (err *ValidationError) GetErrorMessage() string {
	return err.ErrorMessage
}
func (err *ValidationError) SetCode(code string) {
	err.ErrorCode = code
}
func (err *ValidationError) SetError(orgError error) {
	err.OriginError = orgError
}
func (err *ValidationError) SetMessage(msg string) {
	err.ErrorMessage = msg
}

//NotFoundError ..
type NotFoundError errorObj

func (err *NotFoundError) GetCode() string {
	return err.ErrorCode
}
func (err *NotFoundError) GetOriginError() error {
	return err.OriginError
}
func (err *NotFoundError) GetErrorMessage() string {
	return err.ErrorMessage
}
func (err *NotFoundError) SetCode(code string) {
	err.ErrorCode = code
}
func (err *NotFoundError) SetError(orgError error) {
	err.OriginError = orgError
}
func (err *NotFoundError) SetMessage(msg string) {
	err.ErrorMessage = msg
}

//InternalServerError ..
type InternalServerError errorObj

func (err *InternalServerError) GetCode() string {
	return err.ErrorCode
}
func (err *InternalServerError) GetOriginError() error {
	return err.OriginError
}
func (err *InternalServerError) GetErrorMessage() string {
	return err.ErrorMessage
}
func (err *InternalServerError) SetCode(code string) {
	err.ErrorCode = code
}
func (err *InternalServerError) SetError(orgError error) {
	err.OriginError = orgError
}
func (err *InternalServerError) SetMessage(msg string) {
	err.ErrorMessage = msg
}

//UnkownError ..
type UnkownError errorObj

func (err *UnkownError) GetCode() string {
	return err.ErrorCode
}
func (err *UnkownError) GetOriginError() error {
	return err.OriginError
}
func (err *UnkownError) GetErrorMessage() string {
	return err.ErrorMessage
}
func (err *UnkownError) SetCode(code string) {
	err.ErrorCode = code
}
func (err *UnkownError) SetError(orgError error) {
	err.OriginError = orgError
}
func (err *UnkownError) SetMessage(msg string) {
	err.ErrorMessage = msg
}

func GetGrpcError(e Error) error {
	if e != nil {
		var code codes.Code
		switch e.(type) {
		case *ValidationError:
			code = codes.InvalidArgument
		case *NotFoundError:
			code = codes.InvalidArgument
		default:
			code = codes.InvalidArgument
		}
		return status.Errorf(code, e.GetErrorMessage())
	}
	return nil
}
