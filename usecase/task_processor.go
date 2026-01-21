package usecase

import (
	"context"
	"fmt"
	"task-manager-api/repository"
	"time"
)

type TaskProcessor struct {
	taskRepo repository.TaskRepository
}

func NewTaskProcessor(repo repository.TaskRepository) *TaskProcessor {
	return &TaskProcessor{
		taskRepo: repo,
	}
}

// ProcessTaskWithTimeout processes a single task with a timeout
func (p *TaskProcessor) ProcessTaskWithTimeout(ctx context.Context, taskID int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resultCh := make(chan error, 1)

	go func() {
		resultCh <- p.processTask(taskID)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("task processing timeout: %v", ctx.Err())
	case err := <-resultCh:
		return err
	}
}

func (p *TaskProcessor) ProcessTasksInBackground(taskIDs []int) error {
	for _, taskID := range taskIDs {
		go func(id int) {
			p.processTask(id)
		}(taskID)
	}

	return nil
}

func (p *TaskProcessor) processTask(taskID int) error {
	_, err := p.taskRepo.GetByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task %d: %v", taskID, err)
	}

	time.Sleep(2 * time.Second)

	return nil
}
