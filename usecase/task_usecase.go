package usecase

import (
	"context"
	"fmt"
	"task-manager-api/domain"
	"task-manager-api/dto"
	"task-manager-api/repository"
)

type TaskUsecase struct {
	repo  repository.TaskRepository
	cache *CacheService
}

func NewTaskUsecase(repo repository.TaskRepository, cache *CacheService) *TaskUsecase {
	return &TaskUsecase{repo: repo, cache: cache}
}

func (u *TaskUsecase) CreateTask(ctx context.Context, createReq dto.CreateTaskDTO) (dto.TaskResponseDTO, error) {
	if createReq.Title == "" {
		return dto.TaskResponseDTO{}, &domain.ValidationError{Field: "Title", Message: "title is required"}
	}
	if createReq.Description == "" {
		return dto.TaskResponseDTO{}, &domain.ValidationError{Field: "Description", Message: "description is required"}
	}

	task := domain.Task{
		Title:       createReq.Title,
		Description: createReq.Description,
		Status:      createReq.Status,
		Priority:    createReq.Priority,
	}

	var createdTask domain.Task

	err := RetryWithBackoff(ctx, func() error {
		var repoErr error
		createdTask, repoErr = u.repo.Create(task)
		return repoErr
	})

	if err != nil {
		return dto.TaskResponseDTO{}, err
	}

	u.cache.Delete("all_tasks")

	return dto.TaskResponseDTO{
		ID:          createdTask.ID,
		Title:       createdTask.Title,
		Description: createdTask.Description,
		Status:      createdTask.Status,
		Priority:    createdTask.Priority,
	}, nil
}

func (u *TaskUsecase) GetAllTasks() ([]dto.TaskResponseDTO, error) {
	if cached, found := u.cache.Get("all_tasks"); found {
		return cached.([]dto.TaskResponseDTO), nil
	}

	tasks, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responseDTOs []dto.TaskResponseDTO
	for _, task := range tasks {
		responseDTOs = append(responseDTOs, dto.TaskResponseDTO{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
		})
	}

	u.cache.Set("all_tasks", responseDTOs)

	return responseDTOs, nil
}

func (u *TaskUsecase) GetByID(id int) (dto.TaskResponseDTO, error) {
	cacheKey := fmt.Sprintf("task_%d", id)
	if cached, found := u.cache.Get(cacheKey); found {
		return cached.(dto.TaskResponseDTO), nil
	}

	task, err := u.repo.GetByID(id)

	if err != nil {
		return dto.TaskResponseDTO{}, err
	}

	response := dto.TaskResponseDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
	}

	u.cache.Set(cacheKey, response)

	return response, nil
}

func (u *TaskUsecase) UpdateTask(id int, updateReq dto.UpdateTaskDTO) (dto.TaskResponseDTO, error) {
	if updateReq.Title == "" {
		return dto.TaskResponseDTO{}, &domain.ValidationError{Field: "Title", Message: "Title is required"}
	}

	if updateReq.Description == "" {
		return dto.TaskResponseDTO{}, &domain.ValidationError{Field: "Description", Message: "Description is required"}
	}

	existingTask, err := u.repo.GetByID(id)
	if err != nil {
		return dto.TaskResponseDTO{}, err
	}

	existingTask.Title = updateReq.Title
	existingTask.Description = updateReq.Description
	existingTask.Status = updateReq.Status
	existingTask.Priority = updateReq.Priority

	updatedTask, err := u.repo.Update(existingTask)
	if err != nil {
		return dto.TaskResponseDTO{}, err
	}

	cacheKey := fmt.Sprintf("task_%d", id)
	u.cache.Delete(cacheKey)
	u.cache.Delete("all_tasks")

	return dto.TaskResponseDTO{
		ID:          updatedTask.ID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Status:      updatedTask.Status,
		Priority:    updatedTask.Priority,
	}, nil
}

func (u *TaskUsecase) DeleteTask(id int) error {
	_, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}

	err = u.repo.Delete(id)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("task_%d", id)
	u.cache.Delete(cacheKey)
	u.cache.Delete("all_tasks")

	return nil
}
