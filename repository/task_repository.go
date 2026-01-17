package repository

import "task-manager-api/domain"

type TaskRepository interface {
	GetAll() ([]domain.Task, error)
	Create(task domain.Task) (domain.Task, error)
	Close() error
}

type InMemoryTaskRepository struct {
	tasks  []domain.Task
	nextID int
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: []domain.Task{
			{ID: 1, Title: "Task One", Description: "First task description", Status: "Pending", Priority: "High"},
			{ID: 2, Title: "Task Two", Description: "Second task description", Status: "In Progress", Priority: "Medium"},
			{ID: 3, Title: "Task Three", Description: "Third task description", Status: "Completed", Priority: "Low"},
		},
		nextID: 4,
	}
}

func (r *InMemoryTaskRepository) Create(task domain.Task) (domain.Task, error) {
	task.ID = r.nextID
	r.nextID++
	r.tasks = append(r.tasks, task)
	return task, nil
}

func (r *InMemoryTaskRepository) GetAll() ([]domain.Task, error) {
	return r.tasks, nil
}

func (r *InMemoryTaskRepository) Close() error {
	return nil
}
