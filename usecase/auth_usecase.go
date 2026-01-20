package usecase

import (
	"fmt"
	"strings"
	"task-manager-api/domain"
	"task-manager-api/dto"
	"task-manager-api/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *AuthUsecase) Register(req dto.RegisterUserDTO) (dto.UserResponseDTO, error) {
	if !strings.Contains(req.Email, "@") {
		return dto.UserResponseDTO{}, &domain.ValidationError{
			Field:   "email",
			Message: "invalid email format",
		}
	}

	if len(req.Password) < 8 {
		return dto.UserResponseDTO{}, &domain.ValidationError{
			Field:   "password",
			Message: "password must be at least 8 characters long",
		}
	}

	_, err := u.userRepo.GetByEmail(req.Email)
	if err == nil {
		return dto.UserResponseDTO{}, &domain.ValidationError{
			Field:   "email",
			Message: "email already exists",
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponseDTO{}, &domain.DatabaseError{
			Operation: "hash password",
			Err:       err,
		}
	}

	user := domain.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := u.userRepo.Create(user)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	return dto.UserResponseDTO{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}

func (u *AuthUsecase) Login(req dto.LoginUserDTO) (dto.LoginResponseDTO, error) {
	user, err := u.userRepo.GetByEmail(req.Email)
	if err != nil {
		return dto.LoginResponseDTO{}, &domain.AuthenticationError{
			Message: "invalid email or password",
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return dto.LoginResponseDTO{}, &domain.AuthenticationError{
			Message: "invalid email or password",
		}
	}

	token, err := u.generateToken(user.ID)
	if err != nil {
		return dto.LoginResponseDTO{}, &domain.DatabaseError{
			Operation: "generate token",
			Err:       err,
		}
	}

	return dto.LoginResponseDTO{
		Token: token,
		User: dto.UserResponseDTO{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (u *AuthUsecase) ValidateToken(tokenString string) (*domain.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtSecret), nil
	})
	if err != nil {
		return nil, &domain.UnauthorizedError{
			Message: "invalid token",
		}
	}

	if !token.Valid {
		return nil, &domain.UnauthorizedError{
			Message: "invalid token",
		}
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, &domain.UnauthorizedError{
			Message: "invalid token claims",
		}
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, &domain.UnauthorizedError{
			Message: "token expired",
		}
	}

	userID := 0
	_, err = fmt.Sscanf(claims.Subject, "%d", &userID)
	if err != nil {
		return nil, &domain.UnauthorizedError{
			Message: "invalid token subject",
		}
	}

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, &domain.UnauthorizedError{
			Message: "user not found",
		}
	}

	return &user, nil
}

func (u *AuthUsecase) generateToken(userID int) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
