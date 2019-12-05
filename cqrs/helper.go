package cqrs

import (
	"github.com/jedrp/go-core/plresult"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetgRPCError(err plresult.Error) error {
	if err != nil {
		var code codes.Code
		switch err.(type) {
		case *plresult.ValidationError:
			code = codes.InvalidArgument
		case *plresult.NotFoundError:
			code = codes.NotFound
		default:
			code = codes.Internal
		}
		return status.Errorf(code, err.GetErrorMessage())
	}
	return nil
}
