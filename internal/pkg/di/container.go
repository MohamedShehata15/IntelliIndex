package di

import (
	"fmt"
	"sync"
)

// Container manages application dependencies and their lifecycle.
type Container struct {
	mu        sync.RWMutex
	instances map[string]interface{}
	factories map[string]factory
}

type factory func() (interface{}, error)

// NewContainer creates a new dependency injection container.
func NewContainer() *Container {
	return &Container{
		instances: make(map[string]interface{}),
		factories: make(map[string]factory),
	}
}

// Register adds a factory function for a dependency.
func (c *Container) Register(name string, factory factory) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.factories[name] = factory
}

// Resolve gets or creates an instance of the named dependency.
func (c *Container) Resolve(name string) (interface{}, error) {
	if instance, found := c.getExistingInstance(name); found {
		return instance, nil
	}
	return c.createAndStoreInstance(name)
}

// getExistingInstance attempts to retrieve an existing dependency instance.
// Returns the instance and a boolean indicating if it was found.
func (c *Container) getExistingInstance(name string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	instance, ok := c.instances[name]
	return instance, ok
}

// createAndStoreInstance creates a new instance of the dependency and stores it.
func (c *Container) createAndStoreInstance(name string) (interface{}, error) {
	c.mu.Lock()

	if instance, ok := c.instances[name]; ok {
		c.mu.Unlock()
		return instance, nil
	}

	factory, ok := c.factories[name]
	c.mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("no factory registered for dependency: %s", name)
	}

	instance, err := factory()
	if err != nil {
		return nil, fmt.Errorf("error creating dependency %s: %w", name, err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if existingInstance, ok := c.instances[name]; ok {
		return existingInstance, nil
	}

	c.instances[name] = instance
	return instance, nil
}

// MustResolve is like Resolve but panics on error.
func (c *Container) MustResolve(name string) interface{} {
	instance, err := c.Resolve(name)
	if err != nil {
		panic(err)
	}
	return instance
}
