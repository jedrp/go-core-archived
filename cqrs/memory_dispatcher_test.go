package cqrs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/HoaHuynhSoft/go-core/cqrs"
	"github.com/HoaHuynhSoft/go-core/cqrs/mocks"
)

func BenchmarkRegisterHandler(b *testing.B) {
	dispatcher := cqrs.NewInMemoryDispatcher(nil)
	dispatcher.RegisterHandler(context.TODO(), &mocks.IHandler{}, &mocks.ICommand{})
}

func TestCallingHandler(t *testing.T) {
	dispatcher := cqrs.NewInMemoryDispatcher(nil)
	dispatcher.RegisterHandler(context.TODO(), &mocks.IHandler{}, &mocks.ICommand{})

	res := dispatcher.Dispatch(context.TODO(), &mocks.ICommand{})
	fmt.Println(res)
	if !res.IsSuccess {
		t.Errorf("Expect success result but got error one")
	}
}
