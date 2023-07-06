package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/retry"
)

var (
	// to handle the available RPC instances
	instancesLock sync.RWMutex
	instances     = make(map[string][]string) // map to store discovered instances for each service

)

func NewBinaryGenericClient(destServiceName string) (genericclient.Client, error) {
	// **PRE-OPTIMIZATION**
	// instances := DiscoverAddress(destServiceName)

	// **OPTIMIZATION** Service discovery interval: Instead of calling DiscoverAddress every time a new client is made,
	// the instances of valid services are updates every interval as specified in the main() method
	// reason for doing this instead of calling DiscoverAddress every time this method is called is because if there is a sudden surge and large number of requests, NewBinaryGenericCLient is called multiple times which then calls DiscoverAddress multiple times. having it cached can help save unncessary computation in this case

	instances := GetInstances(destServiceName)

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

	lb := loadbalance.NewWeightedRoundRobinBalancer()
	genericCli, err := genericclient.NewClient(destServiceName, generic.BinaryThriftGeneric(), connTimeout, rpcTimeout, client.WithFailureRetry(fp), client.WithLoadBalancer(lb), client.WithHostPorts(instances...))

	return genericCli, err
}
