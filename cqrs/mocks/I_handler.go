// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	"github.com/jedrp/go-core/plresult"
	mock "github.com/stretchr/testify/mock"
)

// IHandler is an autogenerated mock type for the IHandler type
type IHandler struct {
	mock.Mock
}

// Handle provides a mock function with given fields: ctx, command
func (_m *IHandler) Handle(ctx context.Context, command interface{}) *plresult.Result {
	_m.On("Handle", ctx, command).Return(&plresult.Result{IsSuccess: true})
	ret := _m.Called(ctx, command)

	var r0 *plresult.Result
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) *plresult.Result); ok {
		r0 = rf(ctx, command)
	} else {
		r0 = ret.Get(0).(*plresult.Result)
	}

	return r0
}
