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

// HasRole checks if the user has the specified role
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasPermission checks if the user has the specified permission
func (u *User) HasPermission(permission string) bool {
	for _, p := range u.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// AddRole adds a role to the user if they don't already have it
func (u *User) AddRole(role string) {
	if !u.HasRole(role) {
		u.Roles = append(u.Roles, role)
		u.UpdatedAt = time.Now()
	}
}
