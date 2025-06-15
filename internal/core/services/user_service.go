package services

import (
	"errors"
	"time"

	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/incoming"
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"
)

// UserService implements business logic for user operations
type UserService struct {
	securityService   outgoing.SecurityService
	validationService outgoing.ValidationService
}

// NewUserService creates a new user service
func NewUserService(securityService outgoing.SecurityService, validationService outgoing.ValidationService) incoming.UserService {
	return &UserService{
		securityService:   securityService,
		validationService: validationService,
	}
}

// Ensure UserService implements the incoming.UserService interface
var _ incoming.UserService = (*UserService)(nil)

func (u *UserService) CreateUser(username, email, password, displayName string) (*domain.User, error) {
	if err := u.validationService.ValidateEmail(email); err != nil {
		return nil, errors.New("invalid email format")
	}
	if len(username) < 3 {
		return nil, errors.New("username must be at least 3 characters long")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}
	hashedPassword, err := u.securityService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &domain.User{
		ID:          "", // Will be set by repository
		Username:    username,
		Email:       email,
		Password:    hashedPassword,
		Roles:       []string{"user"},
		Permissions: []string{"read:data"},
		DisplayName: displayName,
		APIKeys:     []domain.APIKey{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (u *UserService) ValidateUser(user *domain.User) error {
	if err := u.validationService.ValidateEmail(user.Email); err != nil {
		return errors.New("invalid email format")
	}
	if len(user.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if err := u.validationService.ValidateUsername(user.Username); err != nil {
		return errors.New("invalid username format")
	}
	return nil
}

func (u *UserService) UpdateUserPassword(user *domain.User, password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hashedPassword, err := u.securityService.HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()
	return nil
}

func (u *UserService) CheckUserPassword(user *domain.User, password string) bool {
	return u.securityService.ComparePassword(user.Password, password)
}

func (u *UserService) CreateUserAPIKey(user *domain.User, name string, expiryDays int) (*domain.APIKey, string, error) {
	keyString, err := u.securityService.GenerateSecureToken(32)
	if err != nil {
		return nil, "", err
	}
	hashedKey, err := u.securityService.HashPassword(keyString)
	if err != nil {
		return nil, "", err
	}

	now := time.Now()
	var expiresAt time.Time
	if expiryDays > 0 {
		expiresAt = now.AddDate(0, 0, expiryDays)
	}

	apiKey := domain.APIKey{
		ID:        "", // Will be set by repository
		UserID:    user.ID,
		Name:      name,
		Key:       hashedKey,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}

	user.APIKeys = append(user.APIKeys, apiKey)
	user.UpdatedAt = now

	return &apiKey, keyString, nil
}
