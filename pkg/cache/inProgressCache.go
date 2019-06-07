package cache

import (
	"reflect"
	"sync"
	"time"

	"github.com/box-node-alert-responder/pkg/controller/types"
	log "github.com/sirupsen/logrus"
)


//InProgressCache is a struct to store InProgress of remediation
type InProgressCache struct {
	Items               map[string]types.InProgress
	CacheExpireInterval time.Duration
	Locker              *sync.RWMutex
}

//NewInProgressCache instantiates and returns a new cache
func NewInProgressCache(cacheExpireInterval string) (cache *InProgressCache) {
	interval, _ := time.ParseDuration(cacheExpireInterval)
	return &InProgressCache{
		Items:               make([map[string]types.InProgress),
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
func (cache *InProgressCache) Set(key string, result types.InProgress) bool{
	cache.Locker.Lock()
	defer cache.Locker.Unlock()
	curResult, found := cache.Items[key]
	//No play is currently running
	if !found {
		cache.Items[cond] = result
		return true
	} 
	return false	
}

//GetAll returns current entries of a cache
func (cache *InProgressCache) GetAll() []types.InProgress {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	return cache.Items
}

func (cache *InProgressCache) Count() []types.InProgress {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	return len(cache.Items)
}

//GetItem returns value of a given key and whether it exist or not
func (cache *InProgressCache) GeItem(key string) (types.InProgress, bool) {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	val, found := cache.Items[key]
	if found {
		return val, true
	}
	return types.InProgress{}, false
}

//DelItem deletes a cache item with a given key
func (cache *InProgressCache) DelItem(key string) (types.InProgress, bool) {
	cache.Locker.Lock()
	delete(cache.Items,key)
	cache.Locker.Unlock()
}


