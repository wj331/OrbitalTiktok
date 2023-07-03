package main

import (
	"fmt"
	"log"

	"github.com/nacos-group/nacos-sdk-go/vo"
)

func DiscoverAddress(serviceName string) []string {
	if serviceName == "" {
		return nil
	}

	// **PRE-OPTIMIZATION**
	// cli, err := clients.NewNamingClient(
	// 	vo.NacosClientParam{
	// 		ServerConfigs: []constant.ServerConfig{
	// 			// port of the nacos backend server
	// 			*constant.NewServerConfig("127.0.0.1", 8848),
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to create Nacos client: %v", err)
	// 	return nil
	// }
	// service, err := cli.GetService(vo.GetServiceParam{
	// 	ServiceName: serviceName,
	// })

	// **OPTIMIZATION** initialized the nacos client at the begging instead of creating a new nacos client every time DiscoverAddress is called
	service, err := nacosClient.GetService(vo.GetServiceParam{
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
