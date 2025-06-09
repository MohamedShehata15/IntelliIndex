package di

import "sync"

// Container manages application dependencies and their lifecycle.
type Container struct {
	mu        sync.RWMutex
	instances map[string]interface{}
	factories map[string]factory
}

type factory func() (interface{}, error)
