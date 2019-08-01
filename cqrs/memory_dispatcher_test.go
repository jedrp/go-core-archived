package cqrs_test

import (
	"context"
	"fmt"
	"testing"

	"bitbucket.org/JedHuynhThanh/go-core/cqrs"
	"bitbucket.org/JedHuynhThanh/go-core/cqrs/mocks"
)

func BenchmarkRegisterHandler(b *testing.B) {
	dispatcher := cqrs.NewInMemoryDispatcher()
	dispatcher.RegisterHandler(context.TODO(), &mocks.IHandler{}, &mocks.ICommand{})
}

func TestCallingHandler(t *testing.T) {
	dispatcher := cqrs.NewInMemoryDispatcher()
	dispatcher.RegisterHandler(context.TODO(), &mocks.IHandler{}, &mocks.ICommand{})

	res := dispatcher.Dispatch(context.TODO(), &mocks.ICommand{})
	fmt.Println(res)
	if !res.IsSuccess {
		t.Errorf("Expect success result but got error one")
	}
}
