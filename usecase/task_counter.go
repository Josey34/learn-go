package usecase

import (
	"sync"
	"time"
)

type TaskCounter struct {
	count int
	mu    sync.Mutex
	queue *TaskQueue
}

func NewTaskCounter(capacity int) *TaskCounter {
	return &TaskCounter{count: 0, queue: NewTaskQueue(capacity)}
}

func (tc *TaskCounter) CountTasksWrong(taskIDs []int) int {
	for _, id := range taskIDs {
		tc.mu.Lock()
		tc.count += id
		tc.mu.Unlock()

		go func() {
			tc.mu.Lock()
			tc.count++
			tc.mu.Unlock()
		}()
	}

	time.Sleep(100 * time.Millisecond)
	return tc.count
}

func (tc *TaskCounter) CountTasks(taskIDs []int) (int, error) {
	var wg sync.WaitGroup

	for _, id := range taskIDs {
		err := tc.queue.Enqueue(id)
		if err != nil {
			return 0, err
		}
	}

	for i := 0; i < len(taskIDs); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			id, err := tc.queue.Dequeue()
			if err != nil {
				return
			}

			tc.mu.Lock()
			tc.count += id
			tc.mu.Unlock()
		}()
	}

	wg.Wait()
	return tc.count, nil
}
