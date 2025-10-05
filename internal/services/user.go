package services

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"user-management-go/internal/models"
	"user-management-go/internal/storage"
)

// UserService handles business logic for users
type UserService struct {
	storage   storage.UserStorage
	validator *validator.Validate
}

// NewUserService creates a new user service
func NewUserService(storage storage.UserStorage) *UserService {
	return &UserService{
		storage:   storage,
		validator: validator.New(),
	}
}

// GetAllUsers returns all users
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.storage.GetAll()
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(id string) (*models.User, error) {
	if id == "" {
		return nil, models.ErrUserNotFound
	}
	return s.storage.GetByID(id)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req *models.UserCreateRequest) (*models.User, error) {
	// Validate the request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := s.storage.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, models.ErrDuplicateEmail
	}

	// Generate ID and create user
	id := uuid.New().String()
	user := req.ToUser(id)

	// Additional custom validation
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Save user
	if err := s.storage.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id string, req *models.UserCreateRequest) (*models.User, error) {
	// Check if user exists
	existingUser, err := s.storage.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate the request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Check if email is taken by another user
	if req.Email != existingUser.Email {
		exists, err := s.storage.EmailExists(req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, models.ErrDuplicateEmail
		}
	}

	// Update user data
	updatedUser := req.ToUser(id)
	updatedUser.CreatedAt = existingUser.CreatedAt // Preserve creation time
	updatedUser.UpdatedAt = time.Now()

	// Additional custom validation
	if err := updatedUser.Validate(); err != nil {
		return nil, err
	}

	// Save updated user
	if err := s.storage.Update(updatedUser); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id string) error {
	if id == "" {
		return models.ErrUserNotFound
	}
	return s.storage.Delete(id)
}