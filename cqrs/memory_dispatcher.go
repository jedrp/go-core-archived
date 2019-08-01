package cqrs

import (
	"context"
	"fmt"
	"reflect"

	"github.com/HoaHuynhSoft/go-core/infrastructures"
)

// InMemoryDispatcher ...
type InMemoryDispatcher struct {
	handlers map[string]IHandler
}

// NewInMemoryDispatcher ...
func NewInMemoryDispatcher() *InMemoryDispatcher {
	b := &InMemoryDispatcher{
		handlers: make(map[string]IHandler),
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
func (d *InMemoryDispatcher) Dispatch(ctx context.Context, command interface{}) *infrastructures.Result {
	typeName := reflect.TypeOf(command).String()
	if handler, ok := d.handlers[typeName]; ok {
		return handler.Handle(ctx, command)
	}
	return &infrastructures.Result{IsSuccess: false, Value: "Can not find handler"}
}
