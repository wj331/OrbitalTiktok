package utils

import (
	"fmt"
	"log"
	"sync"

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

func RefreshInstances(services []string) {
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

	// note, if RPC is not currently online, the following line returns the previous available instance from Nacos cache.
	service, err := NacosClient.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
	})
	if err != nil {
		prevInstances := GetInstances(serviceName)
		log.Printf("Failed to discover services: %v", err)
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
