package handler

import (
	"net/http"
	"task-manager-api/usecase"
)

func SetupRoutes(mux *http.ServeMux, uc *usecase.TaskUsecase) {
	RegisterTaskRoutes(mux, uc)
}
