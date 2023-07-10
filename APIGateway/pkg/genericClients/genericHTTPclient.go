package genericClients

import (
	"fmt"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/retry"

	"github.com/simbayippy/OrbitalxTiktok/APIGateway/utils"
)

func newHTTPGenericClient(destServiceName string, thriftFilePath string) (genericclient.Client, error) {
	// **PRE-OPTIMIZATION**
	// instances := DiscoverAddress(destServiceName)

	// **OPTIMIZATION** Service discovery interval: Instead of calling DiscoverAddress every time a new client is made,
	// the instances of valid services are updates every interval as specified in the main() method
	// reason for doing this instead of calling DiscoverAddress every time this method is called is because if there is a sudden surge and large number of requests, NewBinaryGenericCLient is called multiple times which then calls DiscoverAddress multiple times. having it cached can help save unncessary computation in this case
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

	// path := fmt.Sprintf("./thrift/%s.thrift", thriftFilePath)

	p, err := generic.NewThriftFileProvider(thriftFilePath)
	if err != nil {
		klog.Fatalf("new thrift file provider failed: %v", err)
	}
	g, err := generic.HTTPThriftGeneric(p)
	if err != nil {
		klog.Fatalf("new http thrift generic failed: %v", err)
	}

	lb := loadbalance.NewWeightedRoundRobinBalancer()

	// second argument in NewClient() is different from binary client
	genericCli, err := genericclient.NewClient(destServiceName, g, connTimeout, rpcTimeout, client.WithFailureRetry(fp), client.WithLoadBalancer(lb), client.WithHostPorts(instances...))
	if err != nil {
		klog.Fatalf("new http generic client failed: %v", err)
	}

	return genericCli, err
}
