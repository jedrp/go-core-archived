package cqs

import (
	"context"
	"fmt"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/jedrp/go-core/infras"
	"github.com/jedrp/go-core/pllog"
	"google.golang.org/grpc/codes"
)

type MemoryDispatcher struct {
	maxLatencyInMillisecond       time.Duration
	logger                        pllog.PlLogger
	registeredDependencesWrappers map[string]interface{}
}

var (
	INVOKER_INTERNAL_ERROR = infras.Fail(codes.Internal, "An error occurt when server processing the request")
)

func NewMemoryDispatcher(logger pllog.PlLogger, maxLatencyInMillisecond int64) Dispatcher {
	return &MemoryDispatcher{
		maxLatencyInMillisecond:       time.Duration(maxLatencyInMillisecond),
		logger:                        logger,
		registeredDependencesWrappers: make(map[string]interface{}),
	}
}

func (d *MemoryDispatcher) Register(ctx context.Context, deps interface{}, v ...Executor) {
	defer func() {
		if rErr := recover(); rErr != nil {
			d.logger.Panic(rErr, string(debug.Stack()))
		}
	}()
	for _, e := range v {
		typeName := reflect.TypeOf(e).String()
		d.logger.Infof("Registering handler for %s", typeName)
		if _, ok := d.registeredDependencesWrappers[typeName]; ok {
			msg := fmt.Sprintf("Duplicated executer registration detected of type: %s", typeName)
			d.logger.Panic(msg)
		}
		//test
		e.SetDependences(ctx, deps)
		d.registeredDependencesWrappers[typeName] = deps
	}
}

func (d *MemoryDispatcher) Dispatch(ctx context.Context, e Executor) *infras.Result {
	var (
		cancel context.CancelFunc
	)
	if d.maxLatencyInMillisecond > 0 {
		ctx, cancel = context.WithTimeout(ctx, d.maxLatencyInMillisecond*time.Millisecond)
		defer cancel()
	}

	typeName := reflect.TypeOf(e).String()
	if depsWrapper, ok := d.registeredDependencesWrappers[typeName]; ok {
		e.SetDependences(ctx, depsWrapper)
		r := e.Execute(ctx)
		if r.Error != nil {
			pllog.CreateLogEntryFromContext(ctx, d.logger).Error(r.Error.Err())
		}
		return r
	}

	msg := fmt.Sprintf("MemoryDispatcher can't find dependences for type %s", reflect.TypeOf(e).String())
	pllog.CreateLogEntryFromContext(ctx, d.logger).Error(msg)
	return INVOKER_INTERNAL_ERROR
}
