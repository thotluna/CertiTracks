package repositories

import (

	"certitrack/internal/models"
	"errors"

	"time"

	"gorm.io/gorm"
)

var (
	ErrUserExists = errors.New("user with this email already exists")
)

type UserRepository interface {
	CreateUser(user *models.User) error
	EmailExists(email string) bool
	UpdateLastLogin(id string, t time.Time) error
	FindActiveByEmail(email string) (*models.User, error)
	FindActiveByID(id string) (*models.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
    return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(user *models.User) error {

	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) EmailExists(email string) bool {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return false
	}
	return true
}

func (r *UserRepositoryImpl) UpdateLastLogin(id string, t time.Time) error {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}
	user.LastLogin = &t
	return r.db.Save(&user).Error
}

func (r *UserRepositoryImpl) FindActiveByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ? AND is_active = ?", email, true).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindActiveByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ? AND is_active = ?", id, true).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
