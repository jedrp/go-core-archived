package platform

type plError interface {
	GetCode() string
	GetOriginError() error
	GetErrorMessage() string
}

//Error wrapper
type Error struct {
	OriginError  error  //the origin error causing problem
	ErrorCode    string //System specific error code
	ErrorMessage string //System specfic error message
}

func (err *Error) GetCode() string {
	return err.ErrorCode
}
func (err *Error) GetOriginError() error {
	return err.OriginError
}
func (err *Error) GetErrorMessage() string {
	return err.ErrorMessage
}

//ValidationError ..
type ValidationError Error

func (err *ValidationError) GetCode() string {
	return err.ErrorCode
}
func (err *ValidationError) GetOriginError() error {
	return err.OriginError
}
func (err *ValidationError) GetErrorMessage() string {
	return err.ErrorMessage
}

//NotFoundError ..
type NotFoundError Error

func (err *NotFoundError) GetCode() string {
	return err.ErrorCode
}
func (err *NotFoundError) GetOriginError() error {
	return err.OriginError
}
func (err *NotFoundError) GetErrorMessage() string {
	return err.ErrorMessage
}

//InternalServerError ..
type InternalServerError Error

func (err *InternalServerError) GetCode() string {
	return err.ErrorCode
}
func (err *InternalServerError) GetOriginError() error {
	return err.OriginError
}
func (err *InternalServerError) GetErrorMessage() string {
	return err.ErrorMessage
}

//UnkownError ..
type UnkownError Error

func (err *UnkownError) GetCode() string {
	return err.ErrorCode
}
func (err *UnkownError) GetOriginError() error {
	return err.OriginError
}
func (err *UnkownError) GetErrorMessage() string {
	return err.ErrorMessage
}
