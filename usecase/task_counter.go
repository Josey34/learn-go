package usecase

import (
	"sync"
	"time"
)

type TaskCounter struct {
	count int
	mu    sync.Mutex
}

func NewTaskCounter() *TaskCounter {
	return &TaskCounter{count: 0}
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
