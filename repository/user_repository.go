package repository

import (
	"errors"

	"github.com/Arifur999/spotsync/models"

	"gorm.io/gorm"
)

var ErrEmailAlreadyExists = errors.New("user with this email already exists")

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
