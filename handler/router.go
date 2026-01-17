package handler

import "task-manager-api/repository"

func SetupRoutes(repo repository.TaskRepository) {
	RegisterTaskRoutes(repo)
}
