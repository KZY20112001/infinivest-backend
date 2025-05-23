package services

import (
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type ProfileService interface {
	CreateProfile(userID uint, dto dto.ProfileRequest) error
	UpdateProfile(userID uint, dto dto.ProfileRequest) error
	GetProfile(userID uint) (*models.Profile, error)
}

type profileServiceImpl struct {
	repo        repositories.ProfileRepo
	userService UserService
}

func NewProfileServiceImpl(pr repositories.ProfileRepo, us UserService) *profileServiceImpl {
	return &profileServiceImpl{repo: pr, userService: us}
}

func (ps *profileServiceImpl) CreateProfile(userID uint, dto dto.ProfileRequest) error {
	user, err := ps.userService.GetUser(userID)
	if err != nil {
		return err
	}
	profile := models.Profile{
		UserID:            userID,
		User:              *user,
		FirstName:         dto.FirstName,
		LastName:          dto.LastName,
		ProfileUrl:        dto.ProfileUrl,
		ProfileID:         dto.ProfileID,
		RiskTolerance:     dto.RiskTolerance,
		InvestmentStyle:   dto.InvestmentStyle,
		InvestmentHorizon: dto.InvestmentHorizon,
		AnnualIncome:      dto.AnnualIncome,
		ExperienceLevel:   dto.ExperienceLevel,
	}
	return ps.repo.CreateProfile(&profile)
}

func (ps *profileServiceImpl) UpdateProfile(userID uint, dto dto.ProfileRequest) error {
	profile, err := ps.repo.GetProfile(userID)
	if err != nil {
		return err
	}
	profile.FirstName = dto.FirstName
	profile.LastName = dto.LastName
	profile.ProfileUrl = dto.ProfileUrl
	profile.ProfileID = dto.ProfileID
	profile.RiskTolerance = dto.RiskTolerance
	profile.InvestmentStyle = dto.InvestmentStyle
	profile.InvestmentHorizon = dto.InvestmentHorizon
	profile.AnnualIncome = dto.AnnualIncome
	profile.ExperienceLevel = dto.ExperienceLevel
	return ps.repo.UpdateProfile(profile)
}

func (ps *profileServiceImpl) GetProfile(userID uint) (*models.Profile, error) {
	return ps.repo.GetProfile(userID)
}
