package usecase

import "time"

type TaskCounter struct {
	count int
}

func NewTaskCounter() *TaskCounter {
	return &TaskCounter{count: 0}
}

func (tc *TaskCounter) CountTasksWrong(taskIDs []int) int {
	for _, id := range taskIDs {
		tc.count += id
		go func() {
			tc.count++
		}()
	}

	time.Sleep(100 * time.Millisecond)
	return tc.count

}
