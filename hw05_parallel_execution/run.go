package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded   = errors.New("errors limit exceeded")
	ErrInvalidCountConsumers = errors.New("invalid count consumers")
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
		return ErrInvalidCountConsumers
	}

	if m < 1 {
		return ErrErrorsLimitExceeded
	}

	if len(tasks) < 1 {
		return nil
	}

	chanTasksSize := 10
	errCount := ErrCount{
		count: m,
	}

	chanTasks := chanFactory(chanTasksSize)
	produce(chanTasks, &tasks)
	consume(chanTasks, n, &errCount)

	if errCount.Count() < 1 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func consume(chanTasks <-chan Task, n int, errCount *ErrCount) {
	wg := &sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)

		go consumer(chanTasks, errCount, wg)
	}

	wg.Wait()
}

func consumer(chanTasks <-chan Task, errCount *ErrCount, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range chanTasks {
		if errCount.Count() < 1 {
			return
		}

		if err := task(); err != nil {
			errCount.Decrement()
		}
	}
}

func produce(chanTasks chan<- Task, tasks *[]Task) {
	go producer(chanTasks, tasks)
}

func producer(chanTasks chan<- Task, tasks *[]Task) {
	defer close(chanTasks)

	for _, task := range *tasks {
		chanTasks <- task
	}
}

func chanFactory(size int) chan Task {
	chanTasks := make(chan Task, size)

	return chanTasks
}
