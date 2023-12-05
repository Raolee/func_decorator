package v3

import (
	"context"
	"sync"
)

type Stage interface {
	Run(ctx context.Context, input any) ([]any, error)
}

type ConcurrentStage struct {
	Tasks []Task
}

func NewConcurrentStage(tasks ...Task) Stage {
	return &ConcurrentStage{Tasks: tasks}
}

func (s *ConcurrentStage) Run(ctx context.Context, input any) ([]any, error) {
	var wg sync.WaitGroup
	results := make([]any, len(s.Tasks))
	errors := make([]error, len(s.Tasks))

	for i, task := range s.Tasks {
		wg.Add(1)
		go func(i int, t Task) {
			defer wg.Done()
			result, err := t.Execute(ctx, input)
			results[i] = result
			errors[i] = err
		}(i, task)
	}

	wg.Wait()

	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}
