package domain

import "time"

// User represents a user
type User struct {
	ID          string
	Username    string
	Email       string
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

// RemoveRole removes a role from the user
func (u *User) RemoveRole(role string) {
	for i, r := range u.Roles {
		if r == role {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			u.UpdatedAt = time.Now()
			break
		}
	}
}

// AddPermission adds a permission to the user if they don't already have it
func (u *User) AddPermission(permission string) {
	if !u.HasPermission(permission) {
		u.Permissions = append(u.Permissions, permission)
		u.UpdatedAt = time.Now()
	}
}

// RemovePermission removes a permission from the user
func (u *User) RemovePermission(permission string) {
	for i, p := range u.Permissions {
		if p == permission {
			u.Permissions = append(u.Permissions[:i], u.Permissions[i+1:]...)
			u.UpdatedAt = time.Now()
			break
		}
	}
}

// RemoveAPIKey removes an API key from the user
func (u *User) RemoveAPIKey(keyID string) {
	for i, k := range u.APIKeys {
		if k.ID == keyID {
			u.APIKeys = append(u.APIKeys[:i], u.APIKeys[i+1:]...)
			u.UpdatedAt = time.Now()
			break
		}
	}
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	u.LastLoginAt = time.Now()
}
