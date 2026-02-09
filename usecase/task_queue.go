package usecase

import "task-manager-api/domain"

type TaskQueue struct {
	tasks chan int
}

func NewTaskQueue(capacity int) *TaskQueue {
	taskQueue := make(chan int, capacity)
	return &TaskQueue{
		tasks: taskQueue,
	}
}

func (tq *TaskQueue) Enqueue(taskID int) error {
	select {
	case tq.tasks <- taskID:
		return nil
	default:
		return domain.ErrQueueFull
	}
}

func (tq *TaskQueue) Dequeue() (int, error) {
	select {
	case id := <-tq.tasks:
		return id, nil
	default:
		return 0, domain.ErrQueueEmpty
	}
}

func (tq *TaskQueue) Size() int {
	return len(tq.tasks)
}

func (tq *TaskQueue) Close() {
	close(tq.tasks)
}

func (tq *TaskQueue) EnqueueBlocking(taskID int) error {
	tq.tasks <- taskID
	return nil
}
