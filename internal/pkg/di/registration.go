package di

// AdapterRegistrar defines the interface for adapter registration functions
type AdapterRegistrar interface {
	Register(container *Container) error
}

// RegistrationFunc is a function type that implements AdapterRegistrar
type RegistrationFunc func(container *Container) error

// Register implements AdapterRegistrar interface
func (f RegistrationFunc) Register(container *Container) error {
	return f(container)
}

// BatchRegister registers multiple adapters and returns the first error encountered
func BatchRegister(container *Container, registrars ...AdapterRegistrar) error {
	for _, registrar := range registrars {
		if err := registrar.Register(container); err != nil {
			return err
		}
	}
	return nil
}
