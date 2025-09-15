package models

import "errors"

var (
	ErrInvalidPhoneFormat = errors.New("phone must be in format 123-456-7890")
	ErrInvalidFirstName   = errors.New("last name must contain only letters")
	ErrInvalidLastName    = errors.New("last name must contain only letters")
	ErrUserNotFound       = errors.New("user not found")
	ErrDuplicateEmail     = errors.New("email already exisits")
)
