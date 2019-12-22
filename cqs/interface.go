package cqs

import (
	"context"

	"github.com/jedrp/go-core/infras"
)

// Executor interface, runable for Dispatcher
type Executor interface {
	Execute(context.Context) *infras.Result
	SetDependences(context.Context, interface{})
}

// Command interface, change state action
type Command interface {
	Executor
	IsCommand() []string
}

// Query interface, return result, don't change state
type Query interface {
	Executor
	IsQuery() []string
}

// Dispatcher execute command or query, log when command or query return fail status
type Dispatcher interface {
	Dispatch(ctx context.Context, e Executor) *infras.Result
	Register(ctx context.Context, deps interface{}, v ...Executor)
}
