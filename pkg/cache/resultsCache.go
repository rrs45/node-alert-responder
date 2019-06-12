package cache

import (
	"sync"
	"time"

	"github.com/box-node-alert-responder/pkg/types"

	log "github.com/sirupsen/logrus"
)


//ResultsCache is a struct to store results of remediation
type ResultsCache struct {
	Items               map[string]types.ActionResult
	CacheExpireInterval time.Duration
	Locker              *sync.RWMutex
}

//NewResultsCache instantiates and returns a new cache
func NewResultsCache(cacheExpireInterval string) *ResultsCache {
	interval, _ := time.ParseDuration(cacheExpireInterval)
	return &ResultsCache{
		Items:               make(map[string]types.ActionResult),
		CacheExpireInterval: interval,
		Locker:              new(sync.RWMutex),
	}
}

//PurgeExpired expires cache items older than specified purge interval
func (cache *ResultsCache) PurgeExpired() {
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
func (cache *ResultsCache) Set(cond string, t time.Time , worker string, success bool) {
	cache.Locker.Lock()
	var retryCount int
	if prevResult, found := cache.Items[cond]; found {
		log.Infof("Results Cache - %s found in cache", cond)
		if success {
			log.Infof("Results Cache - Current action was successful for %s", cond)
			retryCount = 0
		} else {
				log.Infof("Results Cache - Current action failed for %s", cond)
				retryCount = prevResult.Retry + 1
			}
	} else {
		log.Infof("Results Cache - %s not found in cache", cond)
		if success {
			log.Infof("Results Cache - Current action was successful for %s", cond)
			retryCount = 0
		} else {
			log.Infof("Results Cache - Current action failed for %s", cond)
			retryCount = 1
		}
		
	}

	cache.Items[cond] = types.ActionResult{
		Timestamp: t,
		Worker: worker,
		Success: success,
		Retry: retryCount,
	}
	cache.Locker.Unlock()
}


//GetAll returns current entries of a cache
func (cache *ResultsCache) GetAll() map[string]types.ActionResult {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	return cache.Items
}

//GetItem returns value of a given key and whether it exist or not
func (cache *ResultsCache) GetItem(key string) (types.ActionResult, bool) {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	val, found := cache.Items[key]
	if found {
		return val, true
	}
	return types.ActionResult{}, false
}
