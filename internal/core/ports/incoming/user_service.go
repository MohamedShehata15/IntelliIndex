package incoming

import "github.com/mohamedshehata15/intelli-index/internal/core/domain"

// UserService defines the interface for user business operations
type UserService interface {
	CreateUser(username, email, password, displayName string) (*domain.User, error)
	ValidateUser(user *domain.User) error
	UpdateUserPassword(user *domain.User, password string) error
	CheckUserPassword(user *domain.User, password string) bool
	CreateUserAPIKey(user *domain.User, name string, expiryDays int) (*domain.APIKey, string, error)
}
