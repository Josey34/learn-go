package usecase

import (
	"context"
	"strings"
	"task-manager-api/domain"
	"task-manager-api/repository"
	"time"
)

type TaskSearch struct {
	taskRepo repository.TaskRepository
}

func NewTaskSearch(repo repository.TaskRepository) *TaskSearch {
	return &TaskSearch{taskRepo: repo}
}

func (t *TaskSearch) SearchConcurrently(ctx context.Context, keyword string) ([]domain.Task, error) {
	var titleResults, descResults []domain.Task
	d := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()

	titleResultCh := make(chan []domain.Task, 1)
	descResultCh := make(chan []domain.Task, 1)

	go func() {
		titleResultCh <- t.SearchInTitle(ctx, keyword)
	}()

	go func() {
		descResultCh <- t.SearchInDescription(ctx, keyword)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case titleResults = <-titleResultCh:
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case descResults = <-descResultCh:
	}

	combined := t.mergeResults(titleResults, descResults)
	return combined, nil
}

func (t *TaskSearch) SearchInTitle(ctx context.Context, keyword string) []domain.Task {
	allTasks, err := t.taskRepo.GetAll()

	if err != nil {
		return nil
	}

	var results []domain.Task

	for _, task := range allTasks {
		if strings.Contains(strings.ToLower(task.Title), strings.ToLower(keyword)) {
			results = append(results, task)
		}

	}

	return results
}

func (t *TaskSearch) SearchInDescription(ctx context.Context, keyword string) []domain.Task {
	allTasks, err := t.taskRepo.GetAll()

	if err != nil {
		return nil
	}

	var results []domain.Task

	for _, task := range allTasks {
		if strings.Contains(strings.ToLower(task.Description), strings.ToLower(keyword)) {
			results = append(results, task)
		}

	}

	return results
}

func (t *TaskSearch) mergeResults(titleResults, descResults []domain.Task) []domain.Task {
	seen := make(map[int]bool)

	var merged []domain.Task

	for _, task := range titleResults {
		if !seen[task.ID] {
			seen[task.ID] = true
			merged = append(merged, task)
		}
	}

	for _, task := range descResults {
		if !seen[task.ID] {
			seen[task.ID] = true
			merged = append(merged, task)
		}
	}

	return merged
}
