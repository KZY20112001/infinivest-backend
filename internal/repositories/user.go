package repositories

import (
	"errors"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/models"

	"gorm.io/gorm"
)

type UserRepo interface {
	SignUp(user *models.User) error
	GetUser(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type postgresUserRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepo(db *gorm.DB) *postgresUserRepo {
	return &postgresUserRepo{
		db: db,
	}
}

func (r *postgresUserRepo) SignUp(user *models.User) error {
	if user == nil {
		return commons.ErrNil
	}

	if err := r.db.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return gorm.ErrDuplicatedKey
		}
		return err
	}
	return nil
}

func (r *postgresUserRepo) GetUser(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("ID = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, gorm.ErrDuplicatedKey
		}
		return nil, err
	}

	return &user, nil
}

func (r *postgresUserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("Email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, gorm.ErrDuplicatedKey
		}
		return nil, err
	}

	return &user, nil
}
