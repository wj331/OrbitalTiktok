## Configurations

### File: configs.json allows for the configuration of:
* `APIEndPoint`: the IP & Port of this API Gateway
* `NacosIpAddr` & `NacosPort`: Ip & Port of the Nacos backend server
* `CachingEnabled`: set true to enable caching on all endpoints. Default set to false
* `MaxQPS`: the maximum number of requests per second allowed from a single IP address.
* `BurstSize`: the maximum number of events that can occur at once.

### File: service_configs.json is the configuration file for services. It has the structure (in JSON):
```
{
    "ServiceName": [
      {
        "Version": "v1.0.0",
        "FilePath": "./example/path",
        "GenericClientType": 1
      },
        {
        "Version": "v2.0.0",
        "FilePath": "./example/path2",
        "GenericClientType": 2
      }
    ],
}
```

### GenericClientType integer identifies
* 1: creates a generic Binary client
* 2: creates a generic HTTP client
* 3: creates a generic JSON Thrift client
* 4: creates a generic JSON Protobuf client