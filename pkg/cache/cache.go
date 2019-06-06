package cache

import (
	"sync"
	"time"

	"github.com/box-node-alert-responder/pkg/controller/types"
	log "github.com/sirupsen/logrus"
)

//CacheMap is a struct to store results of remediation
type CacheMap struct {
	Items               map[string]types.ActionResult
	CacheExpireInterval time.Duration
	Locker              *sync.RWMutex
}

//NewCache instantiates and returns a new cache
func NewCache(cacheExpireInterval string) (cache *CacheMap) {
	interval, _ := time.ParseDuration(cacheExpireInterval)
	return &CacheMap{
		Items:               make(map[string]types.ActionResult),
		CacheExpireInterval: interval,
		Locker:              new(sync.RWMutex),
	}
}

//PurgeExpired expires cache items older than specified purge interval
func (cache *CacheMap) PurgeExpired() {
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

//Set creates an entry in the map if it doesnt exist
// or overwrites the timestamp and retry count if it exists
func (cache *CacheMap) Set(node string, issue string, result types.ActionResult) {
	cond := node + "_" + issue
	cache.Locker.Lock()
	curResult, ok := cache.Items[cond]
	if !ok {
		cache.Items[cond] = result
	} else {
		curResult.Timestamp = result.Timestamp
		curResult.Retry++
		cache.Items[cond] = curResult
	}
	cache.Locker.Unlock()
}

//GetAll returns current entries of a cache
func (cache *CacheMap) GetAll() map[string]types.ActionResult {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	return cache.Items
}

//GetItem returns value of a given key and whether it exist or not
func (cache *CacheMap) GetItem(key string) (types.ActionResult, bool) {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	val, found := cache.Items[key]
	if found {
		return val, true
	}
	return types.ActionResult{}, false
}
