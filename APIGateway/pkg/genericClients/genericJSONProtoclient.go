package genericClients

import (
	"fmt"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/retry"

	// "github.com/simbayippy/kitex-repo/pkg/generic"

	"github.com/simbayippy/OrbitalxTiktok/APIGateway/utils"
)

func newJSONProtoGenericClient(destServiceName string, protoFilePath string) (genericclient.Client, error) {
	instances := utils.GetInstances(destServiceName)

	if len(instances) == 0 {
		fmt.Print("No instances found!\n")
		return nil, fmt.Errorf("no instances found for service %s", destServiceName)
	}

	// time before timeout for connecting to RPC server
	connTimeout := client.WithConnectTimeout(100 * time.Millisecond)

	// time before timeout for response from RPC server
	rpcTimeout := client.WithRPCTimeout(3 * time.Second)

	fp := retry.NewFailurePolicy()
	fp.WithMaxRetryTimes(3) // default retries: 2

	p, err := generic.NewPbFileProvider(protoFilePath)
	if err != nil {
		panic(err)
	}
	g, err := generic.JSONProtoGeneric(p)
	if err != nil {
		panic(err)
	}

	lb := loadbalance.NewWeightedRoundRobinBalancer()

	genericCli, err := genericclient.NewClient(destServiceName, g, connTimeout, rpcTimeout, client.WithFailureRetry(fp), client.WithLoadBalancer(lb), client.WithHostPorts(instances...))
	if err != nil {
		panic(err)
	}

	return genericCli, err
}
