package repository

import (
	"task-manager-api/domain"
)

type TaskRepository interface {
	GetAll() ([]domain.Task, error)
	Create(task domain.Task) (domain.Task, error)
	GetByID(id int) (domain.Task, error)
	Update(task domain.Task) (domain.Task, error)
	Delete(id int) error
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

func (r *InMemoryTaskRepository) GetByID(id int) (domain.Task, error) {
	for _, task := range r.tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return domain.Task{}, &domain.NotFoundError{Resource: "Task", ID: id}
}

func (r *InMemoryTaskRepository) Update(updatedTask domain.Task) (domain.Task, error) {
	for i, task := range r.tasks {
		if task.ID == updatedTask.ID {
			r.tasks[i] = updatedTask
			return updatedTask, nil
		}
	}

	return domain.Task{}, &domain.NotFoundError{Resource: "Task", ID: updatedTask.ID}
}

func (r *InMemoryTaskRepository) Delete(id int) error {
	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}

	return &domain.NotFoundError{Resource: "Task", ID: id}
}

func (r *InMemoryTaskRepository) Close() error {
	return nil
}
