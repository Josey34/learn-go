package repository

import "task-manager-api/domain"

type UserRepository interface {
	Create(user domain.User) (domain.User, error)
	GetByEmail(email string) (domain.User, error)
	GetByID(id int) (domain.User, error)
	Close() error
}
