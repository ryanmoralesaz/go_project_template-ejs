package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"user-management-go/components/forms"
	"user-management-go/components/layout"
	"user-management-go/components/pages"
	"user-management-go/internal/models"
	"user-management-go/internal/services"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ShowForm displays the user registration form
func (h *UserHandler) ShowForm(c *gin.Context) {
	// Create the component hierarchy
	formComponent := forms.UserForm()
	component := layout.Base("User Registration", formComponent)

	c.Header("Content-Type", "text/html")
	component.Render(c.Request.Context(), c.Writer)
}

// CreateUser handles form submission to create a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.UserCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to create user"

		// Handle specific errors
		switch err {
		case models.ErrDuplicateEmail:
			statusCode = http.StatusConflict
			message = "Email already exists"
		case models.ErrInvalidPhoneFormat, models.ErrInvalidFirstName, models.ErrInvalidLastName:
			statusCode = http.StatusBadRequest
			message = err.Error()
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":    true,
		"message":    "User created successfully",
		"data":       user,
		"redirectTo": "/users/list",
	})
}

// ListUsers displays all users
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve users",
			"error":   err.Error(),
		})
		return
	}

	// Create the component hierarchy
	usersListComponent := pages.UsersList(users)
	component := layout.Base("Users List", usersListComponent)

	c.Header("Content-Type", "text/html")
	component.Render(c.Request.Context(), c.Writer)
}

// DeleteUser handles user deletion
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User ID is required",
		})
		return
	}

	err := h.userService.DeleteUser(userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to delete user"

		if err == models.ErrUserNotFound {
			statusCode = http.StatusNotFound
			message = "User not found"
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

// GetUser returns a single user by ID (API endpoint)
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to retrieve user"

		if err == models.ErrUserNotFound {
			statusCode = http.StatusNotFound
			message = "User not found"
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}