package utils

import (
	"fmt"
	"log"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var NacosClient naming_client.INamingClient

func DiscoverAddress(serviceName string) []string {
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
