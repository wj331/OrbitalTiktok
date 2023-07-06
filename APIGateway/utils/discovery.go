package utils

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	NacosClient naming_client.INamingClient

	// to handle the available RPC instances
	instancesLock sync.RWMutex
	instances     = make(map[string][]string) // map to store discovered instances for each service
)

func AddInitialInstance(services []string) {
	for _, service := range services {
		instancesLock.Lock()
		instances[service] = discoverAddress(service)
		instancesLock.Unlock()
	}
}

func RefreshInstances(services []string) {
	go func() {
		// TODO: interval of often the API gateway refreshes and gets available services from nacos backend registry
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
				// Refresh the service addresses every minute
				for _, service := range services {
					instancesLock.Lock()
					instances[service] = discoverAddress(service) // DiscoverAddress must be a global function
					instancesLock.Unlock()
				}
			}
		}
	}()
}

func GetInstances(service string) []string {
	instancesLock.RLock()
	instances := instances[service]
	instancesLock.RUnlock()

	return instances
}

func discoverAddress(serviceName string) []string {
	if serviceName == "" {
		return nil
	}

	// **OPTIMIZATION** initialized the nacos client at the begging instead of creating a new nacos client every time DiscoverAddress is called
	service, err := NacosClient.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
	})
	if err != nil {
		log.Fatalf("Failed to discover services: %v", err)
		return nil
	}

	var instances []string
	for _, instance := range service.Hosts {
		address := fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
		instances = append(instances, address)
	}

	fmt.Printf("valid RPC's for %v are: %v\n\n", serviceName, instances)

	return instances
}
