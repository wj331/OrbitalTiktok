package serviceDetails

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/fsnotify/fsnotify"
	"github.com/simbayippy/OrbitalxTiktok/APIGateway/pkg/genericClients"
	"github.com/simbayippy/OrbitalxTiktok/APIGateway/pkg/routes"
)

var (
	services     map[string][]genericClients.ServiceDetails
	serviceNames []string
)

func WatchServiceChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		klog.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	signalCh := make(chan os.Signal, 1)

	// listen for SIGINT signal (Ctrl+C)
	signal.Notify(signalCh, os.Interrupt)

	// Run a goroutine that listens for file system events
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Printf("\nModified %v\n", event.Name)
					// if there is a change in config file -> update the service mapping
					InitServiceMapping()
				}
			case err := <-watcher.Errors:
				klog.Error("error:", err)
			case <-done:
				return
			case <-signalCh:
				close(done)
				return
			}
		}
	}()

	// Watching the service_configs.json file for any changes
	err = watcher.Add("./configs/service_configs.json")
	if err != nil {
		klog.Fatal(err)
	}
	<-done
}

func InitServiceMapping() {
	// Read the config file
	file, err := os.ReadFile("./configs/service_configs.json")
	if err != nil {
		fmt.Printf("Unable to read config file: %v", err)
		return
	}

	if err := json.Unmarshal(file, &services); err != nil {
		klog.Fatalf("Unable to parse JSON file: %v", err)
	}

	fmt.Printf("New service mapping: %v\n", services)

	routes.Pools = genericClients.InitGenericClientPool(services)
}

func GetServiceNames() []string {
	// if serviceNames is empty (uninitialized) or length of services has changed, then a new []string of serviceName is created.
	// Else, return the EXISTING serviceName []string
	if serviceNames == nil || len(serviceNames) != len(services) {
		serviceNames = make([]string, 0, len(services))
		for serviceName := range services {
			serviceNames = append(serviceNames, serviceName)
		}
	}
	return serviceNames
}
