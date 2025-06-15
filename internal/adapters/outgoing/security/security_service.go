package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"
)

// SecurityService implements the outgoing.SecurityService interface
type SecurityService struct {
	bcryptCost int
}

var _ outgoing.SecurityService = (*SecurityService)(nil)

// NewSecurityService creates a new instance of SecurityService
func NewSecurityService() outgoing.SecurityService {
	return &SecurityService{
		bcryptCost: 12,
	}
}

func (s *SecurityService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func (s *SecurityService) ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *SecurityService) GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
