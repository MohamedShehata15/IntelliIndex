package domain

import "time"

// User represents a user
type User struct {
	ID          string
	Username    string
	Password    string
	Roles       []string
	Permissions []string
	DisplayName string
	APIKeys     []APIKey
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastLoginAt time.Time
}

// APIKey represents an API key for accessing the system
type APIKey struct {
	ID        string
	UserID    string
	Name      string
	Key       string
	CreatedAt time.Time
	ExpiresAt time.Time
	LastUsed  time.Time
}
