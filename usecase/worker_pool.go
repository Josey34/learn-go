package usecase

import (
	"sync"
)

type WorkerPool struct {
	numWorkers int
	jobs       chan int
	results    chan error
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan int, 100),
		results:    make(chan error, numWorkers),
	}
}

func (w *WorkerPool) Start(processor func(int) error) {
	for i := 0; i < w.numWorkers; i++ {
		w.wg.Add(1)
		go w.worker(processor)
	}
}

func (w *WorkerPool) worker(processor func(int) error) {
	defer w.wg.Done()

	for taskId := range w.jobs {
		err := processor(taskId)
		w.results <- err
	}
}

func (w *WorkerPool) Submit(taskID int) {
	w.jobs <- taskID
}

func (w *WorkerPool) Wait() []error {
	close(w.jobs)    // Signal workers that no more jobs
	w.wg.Wait()      // Wait for all workers to finish
	close(w.results) // Close results channel

	var errors []error
	for err := range w.results {
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (w *WorkerPool) Stop() {
	close(w.jobs)
	w.wg.Wait()
	close(w.results)
}
