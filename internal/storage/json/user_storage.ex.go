package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"user-management-go/internal/models"
)

// UserStorage implements the storage interface using JSON files
type UserStorage struct {
	filePath string
	mutex    sync.RWMutex
}

// NewUserStorage creates a new JSON-based user storage
func NewUserStorage(filePath string) (*UserStorage, error) {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create empty file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte("[]"), 0644); err != nil {
			return nil, fmt.Errorf("failed to create data file: %w", err)
		}
	}

	return &UserStorage{
		filePath: filePath,
	}, nil
}

// readUsers reads all users from the JSON file
func (s *UserStorage) readUsers() ([]*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read users file: %w", err)
	}

	var users []*models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to parse users data: %w", err)
	}

	return users, nil
}

// writeUsers writes all users to the JSON file
func (s *UserStorage) writeUsers(users []*models.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write users file: %w", err)
	}

	return nil
}

// GetAll returns all users
func (s *UserStorage) GetAll() ([]*models.User, error) {
	return s.readUsers()
}

// GetByID returns a user by their ID
func (s *UserStorage) GetByID(id string) (*models.User, error) {
	users, err := s.readUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, models.ErrUserNotFound
}

// GetByEmail returns a user by their email
func (s *UserStorage) GetByEmail(email string) (*models.User, error) {
	users, err := s.readUsers()
	if err != nil {
		return nil, err
	}

	email = strings.ToLower(strings.TrimSpace(email))
	for _, user := range users {
		if strings.ToLower(user.Email) == email {
			return user, nil
		}
	}

	return nil, models.ErrUserNotFound
}

// Create adds a new user
func (s *UserStorage) Create(user *models.User) error {
	users, err := s.readUsers()
	if err != nil {
		return err
	}

	// Check for duplicate email
	exists, err := s.EmailExists(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return models.ErrDuplicateEmail
	}

	users = append(users, user)
	return s.writeUsers(users)
}

// Update modifies an existing user
func (s *UserStorage) Update(user *models.User) error {
	users, err := s.readUsers()
	if err != nil {
		return err
	}

	for i, u := range users {
		if u.ID == user.ID {
			users[i] = user
			return s.writeUsers(users)
		}
	}

	return models.ErrUserNotFound
}

// Delete removes a user by ID
func (s *UserStorage) Delete(id string) error {
	users, err := s.readUsers()
	if err != nil {
		return err
	}

	for i, user := range users {
		if user.ID == id {
			// Remove user from slice
			users = append(users[:i], users[i+1:]...)
			return s.writeUsers(users)
		}
	}

	return models.ErrUserNotFound
}

// EmailExists checks if an email is already in use
func (s *UserStorage) EmailExists(email string) (bool, error) {
	_, err := s.GetByEmail(email)
	if err == models.ErrUserNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}