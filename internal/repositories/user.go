package repositories

import (
	"errors"

	"github.com/KZY20112001/infinivest-backend/internal/global"
	"github.com/KZY20112001/infinivest-backend/internal/models"

	"gorm.io/gorm"
)

type UserRepo interface {
	SignUp(user *models.User) error
	GetUser(email string) (*models.User, error)
}

type PostgresUserRepo struct {
	db *gorm.DB
}

// Ensure PostgresUserRepo implements the UserRepo interface.
var _ UserRepo = &PostgresUserRepo{}

func NewPostgresUserRepo(db *gorm.DB) *PostgresUserRepo {
	return &PostgresUserRepo{
		db: db,
	}
}

func (ptr *PostgresUserRepo) SignUp(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if err := ptr.db.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.ErrNotFound
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return global.ErrDuplicate
		}
		return err
	}
	return nil
}

func (ptr *PostgresUserRepo) GetUser(email string) (*models.User, error) {
	var user models.User
	if err := ptr.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.ErrNotFound
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, global.ErrDuplicate
		}
		return nil, err
	}

	return &user, nil
}
