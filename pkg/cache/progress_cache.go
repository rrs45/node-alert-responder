package cache

import (
	"sync"
	"time"

	"github.com/box-node-alert-responder/pkg/types"
	log "github.com/sirupsen/logrus"
)

//InProgressCache is a struct to store InProgress of remediation
type InProgressCache struct {
	Items               map[string]types.InProgress
	CacheExpireInterval time.Duration
	Locker              *sync.RWMutex
}

//NewInProgressCache instantiates and returns a new cache
func NewInProgressCache(cacheExpireInterval string) *InProgressCache {
	interval, _ := time.ParseDuration(cacheExpireInterval)
	return &InProgressCache{
		Items:               make(map[string]types.InProgress),
		CacheExpireInterval: interval,
		Locker:              new(sync.RWMutex),
	}
}

//PurgeExpired expires cache items older than specified purge interval
func (cache *InProgressCache) PurgeExpired() {
	ticker := time.NewTicker(cache.CacheExpireInterval)
	for {
		select {
		case <-ticker.C:
			log.Info("CacheManager - Attempting to delete expired entries")
			cache.Locker.Lock()
			for cond, result := range cache.Items {
				if time.Since(result.Timestamp) > cache.CacheExpireInterval {
					log.Info("CacheManager - Deleting expired entry for ", cond)
					delete(cache.Items, cond)
				}
			}
			cache.Locker.Unlock()

		}
	}
}

//Set appends entry to the slice
func (cache *InProgressCache) Set(key string, action types.InProgress) {
	cache.Locker.Lock()
	defer cache.Locker.Unlock()
	cache.Items[key] = action	
}

//GetAll returns current entries of a cache
func (cache *InProgressCache) GetAll() map[string]types.InProgress {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	return cache.Items
}

//Count returns number of items in cache
func (cache *InProgressCache) Count() int {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	return len(cache.Items)
}

//GetItem returns value of a given key and whether it exist or not
func (cache *InProgressCache) GetItem(key string) (types.InProgress, bool) {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	val, found := cache.Items[key]
	if found {
		return val, true
	}
	return types.InProgress{}, false
}

//DelItem deletes a cache item with a given key
func (cache *InProgressCache) DelItem(key string)  {
	cache.Locker.Lock()
	delete(cache.Items,key)
	cache.Locker.Unlock()
}


