package outgoing

// ValidationService defines operations for data validation
type ValidationService interface {
	ValidateEmail(email string) error
	ValidateUsername(username string) error
}
