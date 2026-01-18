package handler

import (
	"task-manager-api/usecase"
)

func SetupRoutes(uc *usecase.TaskUsecase) {
	RegisterTaskRoutes(uc)
}
