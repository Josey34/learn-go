package handler

import (
	"encoding/json"
	"net/http"
	"task-manager-api/domain"
	"task-manager-api/dto"
	"task-manager-api/usecase"
)

func RegisterAuthRoutes(mux *http.ServeMux, uc *usecase.AuthUsecase) {
	mux.HandleFunc("POST /auth/register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r, uc)
	})

	mux.HandleFunc("POST /auth/login", func(w http.ResponseWriter, r *http.Request) {
		login(w, r, uc)
	})

	mux.HandleFunc("GET /auth/me", func(w http.ResponseWriter, r *http.Request) {
		me(w, r, uc)
	})
}

func register(w http.ResponseWriter, r *http.Request, uc *usecase.AuthUsecase) {
	var req dto.RegisterUserDTO

	// Parse request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Call usecase
	user, err := uc.Register(req)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Marshal and return response
	response, err := json.Marshal(user)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func login(w http.ResponseWriter, r *http.Request, uc *usecase.AuthUsecase) {
	var req dto.LoginUserDTO

	// Parse request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Call usecase
	response, err := uc.Login(req)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Marshal and return response
	responseJSON, err := json.Marshal(response)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func me(w http.ResponseWriter, r *http.Request, uc *usecase.AuthUsecase) {
	w.Header().Set("Content-Type", "application/json")

	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*domain.User)
	if !ok {
		HandleError(w, &domain.UnauthorizedError{
			Message: "user not found in context",
		})
		return
	}

	// Convert to DTO
	userDTO := dto.UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	response, err := json.Marshal(userDTO)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.Write(response)
}
