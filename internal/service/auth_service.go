package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/brianchul/airline_booking/internal/config"
	"github.com/brianchul/airline_booking/internal/repository"
	customError "github.com/brianchul/airline_booking/pkg/errors"
	"github.com/brianchul/airline_booking/pkg/jwt"
	"github.com/brianchul/airline_booking/pkg/utils"
)

type AuthService interface {
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtUtil  *jwt.JWT
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtUtil:  jwt.NewJWT(cfg.JWTSecret),
	}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", customError.ErrInvalidCredentials
		}
		return "", err
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return "", customError.ErrInvalidCredentials
	}

	token, err := s.jwtUtil.SignJWT(user.Email, string(user.Tier), time.Now().Add(24*time.Hour))
	if err != nil {
		return "", err
	}

	return token, nil
}
