package utils

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	hertzCache "github.com/hertz-contrib/cache"
	"github.com/hertz-contrib/cache/persist"
)

var (
	CacheHitCount, CacheMissCount int32
)

func CachingDetails(memoryStore *persist.MemoryStore, cacheExpiryTime time.Duration) app.HandlerFunc {
	return hertzCache.NewCacheByRequestURI(
		memoryStore,
		cacheExpiryTime,
		hertzCache.WithOnHitCache(func(ctx context.Context, c *app.RequestContext) {
			atomic.AddInt32(&CacheHitCount, 1)
		}),
		hertzCache.WithOnMissCache(func(ctx context.Context, c *app.RequestContext) {
			atomic.AddInt32(&CacheMissCount, 1)
		}),
	)
}
