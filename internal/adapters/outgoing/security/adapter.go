package security

import "github.com/mohamedshehata15/intelli-index/internal/pkg/di"

// SecurityAdapterFactory implements the outgoing.AdapterRegistrar interface
type SecurityAdapterFactory struct{}

var _ di.AdapterRegistrar = (*SecurityAdapterFactory)(nil)

func NewSecurityAdapterFactory() *SecurityAdapterFactory {
	return &SecurityAdapterFactory{}
}

// Register implements the AdapterRegistrar interface
func (s *SecurityAdapterFactory) Register(container *di.Container) error {
	return RegisterSecurityAdapters(container)
}

// RegisterSecurityAdapters registers all security-related implementations with the DI container
func RegisterSecurityAdapters(container *di.Container) error {
	// Register security service implementation
	container.Register("securityService", func() (interface{}, error) {
		return NewSecurityService(), nil
	})
	return nil
}
