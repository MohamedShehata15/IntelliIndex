package di

// ServiceProvider configures the DI container with all application services and repositories.
type ServiceProvider struct {
	container *Container
}

// NewServiceProvider creates a new service provider with the given container.
func NewServiceProvider(container *Container) *ServiceProvider {
	return &ServiceProvider{
		container,
	}
}

// RegisterServices registers all services and their dependencies.
func (sp *ServiceProvider) RegisterServices() {

}

// GetContainer returns the underlying container.
func (sp *ServiceProvider) GetContainer() *Container {
	return sp.container
}
