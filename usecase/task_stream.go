package usecase

import (
	"fmt"
)

type TaskStream struct {
	tasks   chan int
	results chan string
}

func NewTaskStream() *TaskStream {
	return &TaskStream{
		tasks:   make(chan int),
		results: make(chan string, 10),
	}
}

func (ts *TaskStream) Start() {
	go func() {
		for taskId := range ts.tasks {
			ts.results <- fmt.Sprintf("Processed task %d", taskId)
		}
	}()
}

func (ts *TaskStream) Send(taskId int) {
	ts.tasks <- taskId
}

func (ts *TaskStream) GetResult() string {
	result := <-ts.results
	return result
}

func (ts *TaskStream) Close() {
	close(ts.tasks)
	close(ts.results)
}
