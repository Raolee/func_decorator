package v3

import (
	"context"
)

type TaskType string

const (
	Composite = TaskType("composite")
)

type TaskBuilder interface {
	AddFunction(fn FunctionType) TaskConverterBuilder
	AddLastFunction(fn FunctionType) TaskBuilder
	Build() Task
}

type TaskConverterBuilder interface {
	AttachConverter(fn FunctionType) TaskBuilder
}

type taskBuilder struct {
	taskType TaskType
	fns      []FunctionType
	cvts     []FunctionType
}

func NewTaskBuilder(t TaskType) TaskBuilder {
	return &taskBuilder{
		taskType: t,
		fns:      make([]FunctionType, 0),
		cvts:     make([]FunctionType, 0),
	}
}

func (t *taskBuilder) AddFunction(fn FunctionType) TaskConverterBuilder {
	t.fns = append(t.fns, fn)
	return t
}

func (t *taskBuilder) AddLastFunction(fn FunctionType) TaskBuilder {
	t.fns = append(t.fns, fn)
	return t
}

func (t *taskBuilder) AttachConverter(fn FunctionType) TaskBuilder {
	t.cvts = append(t.cvts, fn)
	return t
}

func (t *taskBuilder) Build() Task {
	var task Task
	switch t.taskType {
	case Composite:
		task = NewCompositeTask()
	default:
		task = NewCompositeTask()
	}

	// fns 가 cvts 보다는 많을 것
	lenFn, lenCvt := len(t.fns), len(t.cvts)

	for i := 0; i < lenFn; i++ {
		if i < lenFn {
			task.AddFunction(t.fns[i])
		}
		if i < lenCvt {
			task.AddFunction(t.cvts[i])
		}
	}
	return task
}

type Task interface {
	Execute(ctx context.Context, input any) (any, error)
	AddFunction(fn FunctionType)
}

type CompositeTask struct {
	Functions []FunctionType
}

func NewCompositeTask() *CompositeTask {
	return &CompositeTask{
		Functions: []FunctionType{},
	}
}

func (t *CompositeTask) AddFunction(fn FunctionType) {
	t.Functions = append(t.Functions, fn)
}

func (t *CompositeTask) Execute(ctx context.Context, input any) (any, error) {
	var currentInput any = input
	var err error

	for _, fn := range t.Functions {
		currentInput, err = fn(ctx, currentInput)
		if err != nil {
			return nil, err
		}
	}

	return currentInput, nil
}
