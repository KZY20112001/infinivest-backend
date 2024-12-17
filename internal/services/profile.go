package services

import (
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type ProfileService struct {
	repo        repositories.ProfileRepo
	userService UserService
}

func NewProfileService(pr repositories.ProfileRepo, us UserService) *ProfileService {
	return &ProfileService{repo: pr, userService: us}
}

func (ps *ProfileService) CreateProfile(userID uint, profileDto dto.ProfileRequest) error {
	user, err := ps.userService.GetUser(userID)
	if err != nil {
		return err
	}
	profile := models.Profile{
		UserID:    userID,
		User:      *user,
		FirstName: profileDto.FirstName,
		LastName:  profileDto.LastName,
	}
	return ps.repo.CreateProfile(&profile)
}

func (ps *ProfileService) UpdateProfile(userID uint, profileDto dto.ProfileRequest) error {
	profile, err := ps.repo.GetProfile(userID)
	if err != nil {
		return err
	}
	profile.FirstName = profileDto.FirstName
	profile.LastName = profileDto.LastName
	return ps.repo.UpdateProfile(profile)
}

func (ps *ProfileService) GetProfile(userID uint) (*models.Profile, error) {
	return ps.repo.GetProfile(userID)
}
