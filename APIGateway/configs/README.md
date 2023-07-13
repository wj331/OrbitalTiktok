## Configurations

`configs.json` allows for the configuration of:

* `APIEndPoint`: the IP & Port of this API Gateway
* `NacosIpAddr` & `NacosPort`: Ip & Port of the Nacos backend server
* `CachingEnabled`: set true to enable caching on all endpoints. Default set to false
* `MaxQPS`: the maximum number of requests per second allowed from a single IP address.
* `BurstSize`: the maximum number of events that can occur at once.