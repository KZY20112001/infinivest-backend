package repositories

import (
	"errors"

	"github.com/KZY20112001/infinivest-backend/internal/constants"
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

func (ptr *postgresUserRepo) SignUp(user *models.User) error {
	if user == nil {
		return constants.ErrNil
	}

	if err := ptr.db.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return constants.ErrDuplicate
		}
		return err
	}
	return nil
}

func (ptr *postgresUserRepo) GetUser(id uint) (*models.User, error) {
	var user models.User
	if err := ptr.db.Where("ID = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrNotFound
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, constants.ErrDuplicate
		}
		return nil, err
	}

	return &user, nil
}

func (ptr *postgresUserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ptr.db.Where("Email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrNotFound
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, constants.ErrDuplicate
		}
		return nil, err
	}

	return &user, nil
}
