package repositories

import (
	"errors"

	"github.com/KZY20112001/infinivest-backend/internal/constants"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"gorm.io/gorm"
)

type ProfileRepo interface {
	CreateProfile(profile *models.Profile) error
	UpdateProfile(updatedProfile *models.Profile) error
	GetProfile(userID uint) (*models.Profile, error)
}

type PostgresProfileRepo struct {
	db *gorm.DB
}

func NewPostgresProfileRepo(db *gorm.DB) *PostgresProfileRepo {
	return &PostgresProfileRepo{db: db}
}

func (ptr *PostgresProfileRepo) CreateProfile(profile *models.Profile) error {
	if profile == nil {
		return errors.New("profile cannot be nil")
	}
	if err := ptr.db.Create(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return constants.ErrDuplicate
		}
		return err
	}
	return nil
}

func (ptr *PostgresProfileRepo) UpdateProfile(updatedProfile *models.Profile) error {
	if err := ptr.db.Save(&updatedProfile).Error; err != nil {
		return err
	}
	return nil
}

func (ptr *PostgresProfileRepo) GetProfile(userID uint) (*models.Profile, error) {
	var profile models.Profile
	if err := ptr.db.Preload("User").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrNotFound
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, constants.ErrDuplicate
		}
		return nil, err
	}
	return &profile, nil
}
