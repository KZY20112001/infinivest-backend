package services

import (
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repositories.UserRepo
}

func NewUserService(ur repositories.UserRepo) *UserService {
	return &UserService{
		repo: ur,
	}
}

func (us *UserService) SignUp(dto dto.SignUpRequest) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := models.User{
		Email:        dto.Email,
		PasswordHash: string(hash),
	}
	if err := us.repo.SignUp(&user); err != nil {
		return err
	}

	return nil
}

func (us *UserService) GetUser(email string) (*models.User, error) {
	return us.repo.GetUser(email)
}
