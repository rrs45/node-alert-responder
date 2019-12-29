package cache

import (
	"sync"
	"time"

	"github.com/box-node-alert-responder/pkg/types"

	log "github.com/sirupsen/logrus"
)


//ResultsCache is a struct to store results of remediation
type ResultsCache struct {
	Items               map[string]map[string]types.ActionResult
	CacheExpireInterval time.Duration
	FailedCountInterval time.Duration
	Locker              *sync.RWMutex
}

//NewResultsCache instantiates and returns a new cache
func NewResultsCache(cacheExpireInterval string, failedCountInterval string) *ResultsCache {
	expireInterval, _ := time.ParseDuration(cacheExpireInterval)
	failedInterval, _ := time.ParseDuration(failedCountInterval)
	return &ResultsCache{
		Items:               make(map[string]map[string]types.ActionResult),
		CacheExpireInterval: expireInterval,
		FailedCountInterval: failedInterval,
		Locker:              new(sync.RWMutex),
	}
}

//PurgeExpired expires cache items older than specified purge interval
func (cache *ResultsCache) PurgeExpired() {
	ticker := time.NewTicker(cache.CacheExpireInterval)
	for {
		select {
		case <-ticker.C:
			log.Debug("CacheManager - Attempting to delete expired entries")
			cache.Locker.Lock()
			for node, actions  := range cache.Items {
				for action, result := range actions {
					if time.Since(result.Timestamp) > cache.CacheExpireInterval {
						//log.Debug("CacheManager - Deleting expired entry for ", )
						delete(cache.Items[node], action)
					}
				}
				if len(cache.Items[node]) == 0 {
					delete(cache.Items, node)
				}
			}
			cache.Locker.Unlock()

		} 
	}
}

//Set creates an entry in the map if it doesnt exist
// or overwrites the timestamp and increments retry count if it exists
func (cache *ResultsCache) Set(node string, action string, result types.ActionResult) {
	cache.Locker.Lock()
	cache.Locker.Unlock()
	var retryCount int
	if prevAction, nodeFound := cache.Items[node]; nodeFound {
		if prevResult, actionFound := prevAction[action]; actionFound {
			log.Debugf("Results Cache - [node:%s, action:%s] found in cache", node, action)
			if result.Success {
				log.Debugf("Results Cache - [node:%s, action:%s] Current action was successful, resetting retry count", node, action)
				retryCount = 0
			} else {
					log.Debugf("Results Cache - [node:%s, action:%s] Current action failed, incrementing retry count", node, action)
					retryCount = prevResult.Retry + 1
				}
		}
	 } else {
		log.Debugf("Results Cache - [node:%s, action:%s] not found in cache", node, action)
		if result.Success {
			log.Debugf("Results Cache - [node:%s, action:%s] Current action was successful, resetting retry count", node, action)
			retryCount = 0
		} else {
			log.Debugf("Results Cache - [node:%s, action:%s] Current action failed, incrementing retry count", node, action)
			retryCount = 1
		}
		
	}

	cache.Items[node][action] = types.ActionResult{
		Timestamp: result.Timestamp,
		Condition: result.Condition,
		Success: result.Success,
		Retry: retryCount,
		Worker: result.Worker,
		}
	
}


//GetAll returns current entries of a cache
func (cache *ResultsCache) GetAll() map[string]map[string]types.ActionResult {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	return cache.Items
}

//GetItem returns value of a given key and whether it exist or not
func (cache *ResultsCache) GetItem(node string, action string) (types.ActionResult, bool) {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	val, found := cache.Items[node][action]
	if found {
		return val, true
	}
	return types.ActionResult{}, false
}

//GetFailedNodeCount returns count of all unique nodes
func (cache *ResultsCache) GetFailedNodeCount() int {
	count := 0  
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	for _, actions := range cache.Items {
		for _, result := range actions {
			if time.Since(result.Timestamp) < cache.FailedCountInterval && !result.Success {
					count++
			}
		}
	}
	return count
}