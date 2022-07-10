package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func doTasks(newTask <-chan Task, ok chan<- bool) {
	ok <- true

	for task := range newTask {
		if err := task(); err == nil {
			ok <- true
		} else {
			ok <- false
		}
	}
}

func sendTasks(newTask chan<- Task, ok <-chan bool, tasks []Task, m int) error {
	errorNumber := 0

	for _, task := range tasks {
		if !<-ok {
			errorNumber++
		}

		if errorNumber >= m {
			return ErrErrorsLimitExceeded
		}

		newTask <- task
	}

	return nil
}

func Run(tasks []Task, n, m int) error {
	newTask := make(chan Task)
	ok := make(chan bool, n)
	wg := sync.WaitGroup{}

	for i := 0; i < n && i < len(tasks); i++ {
		wg.Add(1)
		go func() {
			doTasks(newTask, ok)
			wg.Done()
		}()
	}

	result := sendTasks(newTask, ok, tasks, m)

	close(newTask)
	wg.Wait()
	return result
}
