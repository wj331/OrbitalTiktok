package utils

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type ServiceInstances struct {
	sync.RWMutex
	Instances []string
}

var (
	NacosClient naming_client.INamingClient

	// to handle the available RPC instances
	instances = make(map[string]*ServiceInstances)
)

func AddInitialInstance(services []string) {
	for _, service := range services {
		instances[service] = &ServiceInstances{
			Instances: discoverAddress(service),
		}
	}
}

func RefreshInstances2(services []string) {
	go func() {
		// TODO: interval of often the API gateway refreshes and gets available services from nacos backend registry
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
				for _, service := range services {
					serviceInstances := instances[service]
					serviceInstances.Lock()
					serviceInstances.Instances = discoverAddress(service)
					serviceInstances.Unlock()
				}
			}
		}
	}()
}

func RefreshInstances(services []string) {
	// TODO: interval of often the API gateway refreshes and gets available services from nacos backend registry
	for _, service := range services {
		serviceInstances := instances[service]
		serviceInstances.Lock()
		serviceInstances.Instances = discoverAddress(service)
		serviceInstances.Unlock()
	}

}

func GetInstances(service string) []string {
	serviceInstances := instances[service]
	serviceInstances.RLock()
	serviceInstancesInstances := serviceInstances.Instances
	serviceInstances.RUnlock()

	return serviceInstancesInstances
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
		prevInstances := GetInstances(serviceName)
		log.Printf("Failed to discover services: %v", err)
		// Return the previously known instances (if any) instead of exiting the application
		return prevInstances
	}

	var instances []string
	for _, instance := range service.Hosts {
		address := fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
		instances = append(instances, address)
	}

	log.Printf("\nvalid RPC's for %v are: %v\n", serviceName, instances)

	return instances
}
