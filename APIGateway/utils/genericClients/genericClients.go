package genericClients

import (
	"log"
	"sync"

	"github.com/cloudwego/kitex/client/genericclient"
)

type ServiceDetails struct {
	Version  string
	FilePath string
	// 1: Binary client
	// 2: HTTP client
	// 3: JSON Thrift client
	// 4: JSON Proto client
	GenericClientType int32
}

func InitGenericClientPool(services map[string][]ServiceDetails) map[string]*sync.Pool {
	pools := make(map[string]*sync.Pool)

	// initializing of generic client pools upon apigateway startup
	for serviceName, details := range services {
		// Set up pools for each version of this service
		for _, detail := range details {
			poolKey := serviceName + "_" + detail.Version
			pools[poolKey] = newClientPool(serviceName, detail.FilePath, detail.GenericClientType)
		}
	}
	return pools
}

func newClientPool(serviceName string, filePath string, genericClientType int32) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			var (
				cc  genericclient.Client
				err error
			)
			// TODO: change it to integer values
			switch genericClientType {
			case 1:
				cc, err = newBinaryGenericClient(serviceName)
			case 2:
				cc, err = newHTTPGenericClient(serviceName, filePath)
			case 3:
				cc, err = newJSONGenericClient(serviceName, filePath)
			case 4:
				cc, err = newJSONProtoGenericClient(serviceName, filePath)
			default:
				log.Print("Invalid service name")
				return nil
			}
			if err != nil {
				log.Print("unable to generate new client")
				return nil
			}
			return cc
		},
	}
}
