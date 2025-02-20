package repositories

import (
	"errors"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"gorm.io/gorm"
)

type ProfileRepo interface {
	CreateProfile(profile *models.Profile) error
	UpdateProfile(updatedProfile *models.Profile) error
	GetProfile(userID uint) (*models.Profile, error)
}

type postgresProfileRepo struct {
	db *gorm.DB
}

func NewPostgresProfileRepo(db *gorm.DB) *postgresProfileRepo {
	return &postgresProfileRepo{db: db}
}

func (r *postgresProfileRepo) CreateProfile(profile *models.Profile) error {
	if profile == nil {
		return commons.ErrNil
	}
	if err := r.db.Create(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return gorm.ErrDuplicatedKey
		}
		return err
	}
	return nil
}

func (r *postgresProfileRepo) UpdateProfile(updatedProfile *models.Profile) error {
	if err := r.db.Save(&updatedProfile).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresProfileRepo) GetProfile(userID uint) (*models.Profile, error) {
	var profile models.Profile
	if err := r.db.Preload("User").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, gorm.ErrDuplicatedKey
		}
		return nil, err
	}
	return &profile, nil
}
