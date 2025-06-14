package outgoing

// SecurityService defines operations for security-related tasks
type SecurityService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) bool
	GenerateSecureToken(length int) (string, error)
}
