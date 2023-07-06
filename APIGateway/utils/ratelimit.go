package utils

import (
	"context"
	"log"
	"net"
	"time"

	"golang.org/x/time/rate"

	"github.com/cloudwego/hertz/pkg/app"

	// imported as `cache`
	"github.com/patrickmn/go-cache"
)

var (
	// IP addresses in the cache expires after 5 minutes of no access, and the library by patrickmn automatically cleans up expired items every 6 minutes.
	limiterCache = cache.New(5*time.Minute, 6*time.Minute)

	// TODO: rate limiting numbers
	MaxQPS    = 1000000000 // Each IP address how many QPS
	BurstSize = 1000000000 // number of events that can occur at ONCE. set HIGH for benchmark purposes

	// TODO: cache time allowed before evicted from cache. i.e. how long stored in cache
	cacheExpiryTime = 2 * time.Second

	// cache counters
	cacheHitCount, cacheMissCount int32
)

// Please note that this code has scalability issues. Each instance would have its own cache of rate limiters, and a client could potentially make more requests than allowed by distributing their requests across multiple instances.
// But still ok for now
func RateLimitMiddleware(next func(context.Context, *app.RequestContext)) func(context.Context, *app.RequestContext) {
	return func(ctx context.Context, r *app.RequestContext) {
		clientIP := ""

		if tcpAddr, ok := r.RemoteAddr().(*net.TCPAddr); ok {
			clientIP = tcpAddr.IP.String()
		} else {
			clientIP = "unknown"
		}

		limiter, found := limiterCache.Get(clientIP)
		if !found {
			// If this is the first time we've seen this IP address, create a rate limiter for it.
			limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(MaxQPS)), BurstSize)
			limiterCache.Set(clientIP, limiter, cache.DefaultExpiration)
		}

		// Check to see if a token is available.
		if !limiter.(*rate.Limiter).Allow() {
			log.Printf("Rate limit exceeded for IP %s", clientIP)
			return
		}

		next(ctx, r)
	}
}
