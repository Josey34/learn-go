package usecase

import (
	"errors"
	"task-manager-api/domain"
	"task-manager-api/dto"
	"task-manager-api/repository"
)

type TaskUsecase struct {
	repo repository.TaskRepository
}

func NewTaskUsecase(repo repository.TaskRepository) *TaskUsecase {
	return &TaskUsecase{repo: repo}
}

func (u *TaskUsecase) CreateTask(createReq dto.CreateTaskDTO) (dto.TaskResponseDTO, error) {
	if createReq.Title == "" {
		return dto.TaskResponseDTO{}, errors.New("title is required")
	}
	if createReq.Description == "" {
		return dto.TaskResponseDTO{}, errors.New("description is required")
	}

	task := domain.Task{
		Title:       createReq.Title,
		Description: createReq.Description,
		Status:      createReq.Status,
		Priority:    createReq.Priority,
	}

	createdTask, err := u.repo.Create(task)

	if err != nil {
		return dto.TaskResponseDTO{}, err
	}

	return dto.TaskResponseDTO{
		ID:          createdTask.ID,
		Title:       createdTask.Title,
		Description: createdTask.Description,
		Status:      createdTask.Status,
		Priority:    createdTask.Priority,
	}, nil
}

func (u *TaskUsecase) GetAllTasks() ([]dto.TaskResponseDTO, error) {
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

	return responseDTOs, nil
}

func (u *TaskUsecase) GetByID(id int) (dto.TaskResponseDTO, error) {
	task, err := u.repo.GetByID(id)

	if err != nil {
		return dto.TaskResponseDTO{}, err
	}

	return dto.TaskResponseDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
	}, nil
}

func (u *TaskUsecase) UpdateTask(id int, updateReq dto.UpdateTaskDTO) (dto.TaskResponseDTO, error) {
	if updateReq.Title == "" {
		return dto.TaskResponseDTO{}, errors.New("Title is required")
	}

	if updateReq.Description == "" {
		return dto.TaskResponseDTO{}, errors.New("Description is required")
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

	return u.repo.Delete(id)
}
