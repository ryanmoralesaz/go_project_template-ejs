package models

import (
	"strings"
	"time"
	"regexp"
)

type User struct {
	ID        string    `json:"id" validate:"required"`
	FirstName string    `json:"firstName" validate:"required,min=2,max=50,alpha"`
	LastName  string    `json:"lastName" validate:"required,min=2,max=50,alpha"`
	Email     string    `json:"email" validate:"required,email"`
	Phone     string    `json:"phone" validate:"required,phone"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserCreateRequest struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=50,alpha"`
	LastName  string `json:"lastName" validate:"required,min=2,max=50,alpha"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,phone"`
}

func (ucr *UserCreateRequest) ToUser(id string) *User {
	now := time.Now()
	return &User{
		ID:        id,
		FirstName: strings.TrimSpace(ucr.FirstName),
		LastName:  strings.TrimSpace(ucr.LastName),
		Email:     strings.TrimSpace(strings.ToLower(ucr.Email)),
		Phone:     strings.TrimSpace(ucr.Phone),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) Validate() error {
	phoneRegex := regexp.MustCompile(`^\d{3}-\d{3}-\d{4}$`)
	if !phoneRegex.MatchString(u.Phone) {
		return ErrInvalidPhoneFormat
	}

	nameRegex := regexp.MustCompile(`^[A-Za-z]+$`)
	if !nameRegex.MatchString(u.FirstName) {
		return ErrInvalidFirstName
	}
	if !nameRegex.MatchString(u.LastName) {
		return ErrInvalidLastName
	}
	return nil
}
