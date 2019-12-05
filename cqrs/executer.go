package cqrs

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/jedrp/go-core/pllog"
	"github.com/jedrp/go-core/plresult"
)

var (
	INVOKER_INTERNAL_ERROR = &plresult.InternalServerError{
		ErrorCode:    "INTERNAL_SERVER_ERROR",
		ErrorMessage: "An error occurt when server processing the request",
		OriginError:  errors.New("An error occurt when server processing the request"),
	}
)

type Executer interface {
	Execute(context.Context)
	GetError() plresult.Error
	SetError(plresult.Error)
	SetDependencesWrapper(context.Context, interface{}) error
}

type Commander interface {
	Executer
}

type Querer interface {
	Executer
}

type Invoker interface {
	RegisterExecuter(context.Context, interface{}, ...Executer) error
	Invoke(context.Context, Executer)
}

type MemoryExecutableInvoker struct {
	maxLatencyInMillisecond       time.Duration
	logger                        pllog.PlLogger
	registeredDependencesWrappers map[string]interface{}
}

func NewMemoryExecutableInvoker(logger pllog.PlLogger, maxLatencyInMillisecond int64) *MemoryExecutableInvoker {
	return &MemoryExecutableInvoker{
		maxLatencyInMillisecond:       time.Duration(maxLatencyInMillisecond),
		logger:                        logger,
		registeredDependencesWrappers: make(map[string]interface{}),
	}
}

func (invoker *MemoryExecutableInvoker) RegisterExecuter(ctx context.Context, depsWrapper interface{}, executers ...Executer) error {
	defer func() {
		if rErr := recover(); rErr != nil {
			invoker.logger.Panic(rErr, string(debug.Stack()))
		}
	}()
	for _, e := range executers {
		typeName := reflect.TypeOf(e).String()
		invoker.logger.Infof("Registering handler for %s", typeName)
		if _, ok := invoker.registeredDependencesWrappers[typeName]; ok {
			msg := fmt.Sprintf("Duplicated executer registration detected of type: %s", typeName)
			invoker.logger.Panic(msg)
		}
		//test
		e.SetDependencesWrapper(ctx, depsWrapper)
		invoker.registeredDependencesWrappers[typeName] = depsWrapper
	}
	return nil
}

func (invoker *MemoryExecutableInvoker) Invoke(ctx context.Context, e Executer) {
	ctx, cancel := context.WithTimeout(ctx, invoker.maxLatencyInMillisecond*time.Millisecond)
	defer cancel()
	typeName := reflect.TypeOf(e).String()
	if depsWrapper, ok := invoker.registeredDependencesWrappers[typeName]; ok {
		e.SetDependencesWrapper(ctx, depsWrapper)
		e.Execute(ctx)
		err := e.GetError()
		if err != nil {
			pllog.CreateLogEntryFromContext(ctx, invoker.logger).Error(err.GetErrorMessage())
		}
		return
	}

	msg := fmt.Sprintf("MemoryExecutableInvoker can't find dependences for type %s", reflect.TypeOf(e).String())
	pllog.CreateLogEntryFromContext(ctx, invoker.logger).Errorf(msg)
	e.SetError(INVOKER_INTERNAL_ERROR)
}
