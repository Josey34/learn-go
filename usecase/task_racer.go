package usecase

import (
	"fmt"
	"math/rand"
	"time"
)

const taskCompletionFormat = "Task %s completed in %d"

type TaskRacer struct {
	task1Chan chan string
	task2Chan chan string
	task3Chan chan string
}

func NewTaskRacer() *TaskRacer {
	return &TaskRacer{
		task1Chan: make(chan string, 10),
		task2Chan: make(chan string, 10),
		task3Chan: make(chan string, 10),
	}
}

func (tr *TaskRacer) StartRace() string {
	go func() {
		randomMs := rand.Intn(300)
		time.Sleep(time.Duration(randomMs) * time.Millisecond)
		tr.task1Chan <- fmt.Sprintf("Task 1 completed in %dms", randomMs)
	}()

	go func() {
		randomMs := rand.Intn(300)
		time.Sleep(time.Duration(randomMs) * time.Millisecond)
		tr.task2Chan <- fmt.Sprintf("Task 2 completed in %dms", randomMs)
	}()

	go func() {
		randomMs := rand.Intn(300)
		time.Sleep(time.Duration(randomMs) * time.Millisecond)
		tr.task3Chan <- fmt.Sprintf("Task 3 completed in %dms", randomMs)
	}()

	select {
	case result := <-tr.task1Chan:
		return result
	case result := <-tr.task2Chan:
		return result
	case result := <-tr.task3Chan:
		return result
	}
}
