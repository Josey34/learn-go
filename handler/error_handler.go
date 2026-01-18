package handler

import (
	"encoding/json"
	"net/http"
	"task-manager-api/domain"
	"task-manager-api/dto"
)

func HandleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch e := err.(type) {
	case *domain.ValidationError:
		// 400 Bad Request
		response := dto.ErrorResponse{
			Error:   "ValidationError",
			Message: e.Error(),
			Details: map[string]string{"field": e.Field},
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	case *domain.NotFoundError:
		// 404 Not Found
		response := dto.ErrorResponse{
			Error:   "NotFound",
			Message: e.Error(),
		}

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)

	case *domain.DatabaseError:
		// 500 Internal Server Error
		response := dto.ErrorResponse{
			Error:   "DatabaseError",
			Message: e.Error(),
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)

	default:
		// 500 Internal Server Error for unknown errors
		response := dto.ErrorResponse{
			Error:   "InternalServerError",
			Message: "An unexpected error occurred",
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
	}
}
