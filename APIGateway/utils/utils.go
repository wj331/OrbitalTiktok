// utils/utils.go
package utils

import (
	"sync"
	"time"
)

type InstanceManager struct {
	lock      sync.RWMutex
	instances map[string][]string
}

func NewInstanceManager() *InstanceManager {
	return &InstanceManager{
		instances: make(map[string][]string),
	}
}

func (im *InstanceManager) AddInitialInstance(services []string) {
	for _, service := range services {
		im.lock.Lock()
		im.instances[service] = DiscoverAddress(service)
		im.lock.Unlock()
	}
}

func (im *InstanceManager) RefreshInstances(services []string) {
	go func() {
		// TODO: interval of often the API gateway refreshes and gets available services from nacos backend registry
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
				// Refresh the service addresses every minute
				for _, service := range services {
					im.lock.Lock()
					im.instances[service] = DiscoverAddress(service) // DiscoverAddress must be a global function
					im.lock.Unlock()
				}
			}
		}
	}()
}

func (im *InstanceManager) GetInstances(service string) []string {
	im.lock.RLock()
	instances := im.instances[service]
	im.lock.RUnlock()

	return instances
}

// more methods to manipulate instances
