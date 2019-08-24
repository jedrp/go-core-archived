package cqrs

import (
	"context"

	"github.com/HoaHuynhSoft/go-core/plresult"
)

// IHandler ...
type IHandler interface {
	Handle(ctx context.Context, command interface{}) *plresult.Result
}

// IDispatcher ...
type IDispatcher interface {
	Dispatch(ctx context.Context, command interface{}) *plresult.Result
	RegisterHandler(ctx context.Context, handler IHandler, commands ...interface{}) error
}

// ICommand ...
type ICommand interface {
	New(interface{}) error
}

// IQuery ...
type IQuery interface {
}
