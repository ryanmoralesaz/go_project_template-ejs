package storage

import "user-database-go/internal/models"

type UserStorage interface {
	GetAll() ([]*models.User, error)
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id string) error
	EmailExists(email string) (bool, error)
}
