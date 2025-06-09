package di

import (
	"log"
	"sync"
)

var (
	defaultContainer *Container
	containerOnce    sync.Once
)

// Bootstrap initializes the default container with all services.
// This function should be called once during application startup.
func Bootstrap() *Container {
	containerOnce.Do(func() {
		defaultContainer = NewContainer()
		provider := NewServiceProvider(defaultContainer)
		provider.RegisterServices()
		log.Println("Dependency injection container initialized")
	})
	return defaultContainer
}

// ResetContainer clears all instances from the container but keeps factories.
func ResetContainer() {
	if defaultContainer != nil {
		defaultContainer.Reset()
	}
}
