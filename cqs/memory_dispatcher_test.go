package cqs_test

import (
	"context"
	"testing"

	"github.com/jedrp/go-core/cqs"
	"github.com/jedrp/go-core/infras"
	"github.com/jedrp/go-core/pllog"
	"google.golang.org/grpc/codes"
)

type testCommand struct{}
type testQuery struct{}
type testDeps struct{}

func (*testQuery) Execute(context.Context) *infras.Result {
	return infras.Fail(codes.Internal, "test internal")
}

func (*testQuery) SetDependences(context.Context, interface{}) {

}

func (*testCommand) Execute(context.Context) *infras.Result {
	return infras.OK(1)
}

func (*testCommand) SetDependences(context.Context, interface{}) {

}
func TestDitpatcher(t *testing.T) {
	d := cqs.NewMemoryDispatcher(
		&pllog.DefaultLogger{},
		100,
	)
	ctx := context.Background()
	d.Register(ctx, &testDeps{}, &testCommand{})
	d.Register(ctx, &testDeps{}, &testQuery{})

	r := d.Dispatch(ctx, &testCommand{})

	if r.Error != nil {
		t.Error("should not return error")
	}
	if r.Value.(int) != 1 {
		t.Errorf("expected 1 but got %v", r.Value)
	}

	r = d.Dispatch(ctx, &testQuery{})

	if r.Error == nil {
		t.Error("should return error")
	}
	if r.Error.Code() != codes.Internal {
		t.Errorf("expected Internal but got %v", r.Error.Code())
	}
}
