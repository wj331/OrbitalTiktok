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

func NewJSONProtoGenericClient(destServiceName string, protoFilePath string) (genericclient.Client, error) {
	instances := utils.GetInstances(destServiceName)

	if len(instances) == 0 {
		fmt.Print("No instances found!\n")
		return nil, fmt.Errorf("no instances found for service %s", destServiceName)
	}

	// time before timeout for connecting to RPC server
	connTimeout := client.WithConnectTimeout(100 * time.Millisecond)

	// time before timeout for response from RPC server
	rpcTimeout := client.WithRPCTimeout(3 * time.Second)

	// import "github.com/cloudwego/kitex/pkg/retry"
	fp := retry.NewFailurePolicy()
	fp.WithMaxRetryTimes(3) // set the maximum number of retries to 3, default 2: client.WithFailureRetry(fp)

	p, err := generic.NewPbFileProvider(protoFilePath)
	if err != nil {
		panic(err)
	}
	// EXTRA ASSIGNMENT PART: under kitex/pkg/generic, calling the method generci.JSONThriftGeneric() also internally calls from the jsonthrift_codec!!
	g, err := generic.JSONProtoGeneric(p)
	if err != nil {
		panic(err)
	}

	lb := loadbalance.NewWeightedRoundRobinBalancer()

	// second argument in NewClient() is different from binary client
	genericCli, err := genericclient.NewClient(destServiceName, g, connTimeout, rpcTimeout, client.WithFailureRetry(fp), client.WithLoadBalancer(lb), client.WithHostPorts(instances...))
	if err != nil {
		panic(err)
	}

	return genericCli, err
}
