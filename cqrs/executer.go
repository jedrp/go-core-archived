package cqrs

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/jedrp/go-core/pllog"
	"github.com/jedrp/go-core/plresult"
)

type Executer interface {
	Execute(context.Context) *plresult.Result
	SetDependencesWrapper(context.Context, interface{}) error
}

type Invoker interface {
	RegisterHandler(context.Context, interface{})
	Invoke(context.Context, Executer) *plresult.Result
}

type MemoryExecutableInvoker struct {
	logger                        pllog.PlLogger
	registeredDependencesWrappers map[string]interface{}
}

func (invoker *MemoryExecutableInvoker) RegisterHandler(ctx context.Context, depsWrapper interface{}, executers ...Executer) error {
	for _, e := range executers {
		typeName := reflect.TypeOf(e).String()
		invoker.logger.Infof("Registering handler for", typeName)
		if _, ok := invoker.registeredDependencesWrappers[typeName]; ok {
			msg := fmt.Sprintf("duplicate command handler registration with command bus for command of type: %s", typeName)
			invoker.logger.Panic(msg)
		}
		invoker.registeredDependencesWrappers[typeName] = depsWrapper
	}
	return nil
}

func (invoker *MemoryExecutableInvoker) Invoke(ctx context.Context, exer Executer) *plresult.Result {
	typeName := reflect.TypeOf(exer).String()
	if depsWrapper, ok := invoker.registeredDependencesWrappers[typeName]; ok {
		err := exer.SetDependencesWrapper(ctx, depsWrapper)
		if err != nil {
			return plresult.InternalErrorResult(fmt.Errorf("MemoryExecutableInvoker can't set handler for type %s", typeName), "FAIL_TO_SET_HANDLER")
		}
		result := exer.Execute(ctx)
		if invoker.logger != nil && !result.IsSuccess && result.Error != nil {
			pllog.CreateLogEntryFromContext(ctx, invoker.logger).Error(result.Error.GetOriginError().Error())
		}
		return result
	}
	return plresult.InternalErrorResult(errors.New("MemoryExecutableInvoker can't find handler"), "MISING_HANDLER_CONFIGURATION")

}
