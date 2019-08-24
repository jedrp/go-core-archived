package cqrs

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/HoaHuynhSoft/go-core/pllog"
	"github.com/HoaHuynhSoft/go-core/plresult"
)

// InMemoryDispatcher ...
type InMemoryDispatcher struct {
	handlers map[string]IHandler
	logger   pllog.PlLogger
}

// NewInMemoryDispatcher ...
func NewInMemoryDispatcher(loggerImpl pllog.PlLogger) *InMemoryDispatcher {
	b := &InMemoryDispatcher{
		handlers: make(map[string]IHandler),
		logger:   loggerImpl,
	}
	return b
}

// RegisterHandler ...
func (d *InMemoryDispatcher) RegisterHandler(ctx context.Context, handler IHandler, commands ...interface{}) error {
	for _, command := range commands {
		typeName := reflect.TypeOf(command).String()
		fmt.Println("Registering handler for", typeName)
		if _, ok := d.handlers[typeName]; ok {
			return fmt.Errorf("duplicate command handler registration with command bus for command of type: %s", typeName)
		}
		d.handlers[typeName] = handler
	}
	return nil
}

// Dispatch ...
func (d *InMemoryDispatcher) Dispatch(ctx context.Context, command interface{}) *plresult.Result {
	typeName := reflect.TypeOf(command).String()
	if handler, ok := d.handlers[typeName]; ok {
		result := handler.Handle(ctx, command)
		if d.logger != nil && !result.IsSuccess {
			d.logger.Error(result.Error.GetOriginError())
		}
		return result
	}
	return plresult.InternalErrorResult(errors.New("InMemoryDispatcher can't find handler"), "MISING_HANDLER_IMPL")
}
