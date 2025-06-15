package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"
)

// ValidationService provides concrete implementation of outgoing.ValidationService
type ValidationService struct {
	emailRegex    *regexp.Regexp
	usernameRegex *regexp.Regexp
}

// NewValidationService creates a new validation service
func NewValidationService() outgoing.ValidationService {
	return &ValidationService{
		emailRegex:    regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		usernameRegex: regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`),
	}
}

var _ outgoing.ValidationService = (*ValidationService)(nil)

// ValidateEmail validates email format and sanitizes it
func (v *ValidationService) ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	sanitizedEmail := strings.ToLower(strings.TrimSpace(email))

	if !v.emailRegex.MatchString(sanitizedEmail) {
		return fmt.Errorf("invalid email format")
	}

	if len(sanitizedEmail) > 254 {
		return fmt.Errorf("email too long")
	}

	return nil
}

// ValidateUsername validates username format and sanitizes it
func (v *ValidationService) ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	sanitizedUsername := strings.TrimSpace(username)

	if !v.usernameRegex.MatchString(sanitizedUsername) {
		return fmt.Errorf("username must be 3-30 characters and contain only letters, numbers, underscores, and hyphens")
	}

	reserved := []string{"admin", "root", "system", "api", "www", "mail", "ftp"}
	for _, r := range reserved {
		if strings.EqualFold(sanitizedUsername, r) {
			return fmt.Errorf("username '%s' is reserved", sanitizedUsername)
		}
	}

	return nil
}
