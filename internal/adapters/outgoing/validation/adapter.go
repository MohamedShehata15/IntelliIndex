package validation

import "github.com/mohamedshehata15/intelli-index/internal/pkg/di"

// ValidationAdapterFactory implements the outgoing.AdapterRegistrar interface
type ValidationAdapterFactory struct{}

var _ di.AdapterRegistrar = (*ValidationAdapterFactory)(nil)

// NewValidationAdapterFactory creates a new factory for validation adapters
func NewValidationAdapterFactory() *ValidationAdapterFactory {
	return &ValidationAdapterFactory{}
}

func (v *ValidationAdapterFactory) Register(container *di.Container) error {
	return RegisterValidationAdapters(container)
}

// RegisterValidationAdapters registers all validation-related implementations with the DI container
func RegisterValidationAdapters(container *di.Container) error {
	// Register validation service implementation
	container.Register("validationService", func() (interface{}, error) {
		return NewValidationService(), nil
	})
	return nil
}
