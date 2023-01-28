package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidJobs         = errors.New("invalid count jobs")
)

type Task func() error

type ErrCount struct {
	mu    sync.RWMutex
	count int
}

func (err *ErrCount) Count() int {
	err.mu.RLock()
	defer err.mu.RUnlock()

	return err.count
}

func (err *ErrCount) Decrement() {
	err.mu.Lock()
	defer err.mu.Unlock()

	err.count--
}

func Run(tasks []Task, n, m int) error {
	if n < 1 {
		return ErrInvalidJobs
	}

	if m < 1 {
		return ErrErrorsLimitExceeded
	}

	if len(tasks) < 1 {
		return nil
	}

	wg := sync.WaitGroup{}
	errCount := ErrCount{
		count: m,
	}

	ch := make(chan Task, len(tasks))

	for _, task := range tasks {
		ch <- task
	}

	close(ch)

	job := func(ch <-chan Task) {
		defer wg.Done()

		for task := range ch {
			if errCount.Count() < 1 {
				return
			}

			if err := task(); err != nil {
				errCount.Decrement()
			}
		}
	}

	wg.Add(n)

	for i := 0; i < n; i++ {
		go job(ch)
	}

	wg.Wait()

	if errCount.Count() < 1 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
