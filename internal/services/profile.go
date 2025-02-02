package services

import (
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type ProfileService interface {
	CreateProfile(userID uint, profileDto dto.ProfileRequest) error
	UpdateProfile(userID uint, profileDto dto.ProfileRequest) error
	GetProfile(userID uint) (*models.Profile, error)
}

type profileServiceImpl struct {
	repo        repositories.ProfileRepo
	userService UserService
}

func NewProfileServiceImpl(pr repositories.ProfileRepo, us UserService) *profileServiceImpl {
	return &profileServiceImpl{repo: pr, userService: us}
}

func (ps *profileServiceImpl) CreateProfile(userID uint, profileDto dto.ProfileRequest) error {
	user, err := ps.userService.GetUser(userID)
	if err != nil {
		return err
	}
	profile := models.Profile{
		UserID:     userID,
		User:       *user,
		FirstName:  profileDto.FirstName,
		LastName:   profileDto.LastName,
		Address:    profileDto.Address,
		ProfileUrl: profileDto.ProfileUrl,
		ProfileID:  profileDto.ProfileID,
	}
	return ps.repo.CreateProfile(&profile)
}

func (ps *profileServiceImpl) UpdateProfile(userID uint, profileDto dto.ProfileRequest) error {
	profile, err := ps.repo.GetProfile(userID)
	if err != nil {
		return err
	}
	profile.FirstName = profileDto.FirstName
	profile.LastName = profileDto.LastName
	profile.Address = profileDto.Address
	profile.ProfileUrl = profileDto.ProfileUrl
	profile.ProfileID = profileDto.ProfileID
	return ps.repo.UpdateProfile(profile)
}

func (ps *profileServiceImpl) GetProfile(userID uint) (*models.Profile, error) {
	return ps.repo.GetProfile(userID)
}
