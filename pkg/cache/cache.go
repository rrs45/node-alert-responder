package cache

import (
	"sync"
	"time"

	_ "github.com/box-node-alert-responder/pkg/controller/types"
	log "github.com/sirupsen/logrus"
)

//ActionResult is a struct to represent result of a remediation
type ActionResult struct {
	Timestamp  time.Time // Time when the play kicked off
	ActionName string    // Script or Ansible play name
	Success    bool      // Whether it was fixed or not
	Retry      int       // Number of times to retry the play if not successful
}

//CacheMap is a struct to store results of remediation
type CacheMap struct {
	Items         map[string]ActionResult
	PurgeInterval time.Duration
	Locker        *sync.RWMutex
}

//NewCache instantiates and returns a new cache
func NewCache(purgeInterval string) (cache *CacheMap) {
	interval, _ := time.ParseDuration(purgeInterval)
	return &CacheMap{
		Items:         make(map[string]ActionResult),
		PurgeInterval: interval,
		Locker:        new(sync.RWMutex),
	}
}

//PurgeExpired expires cache items older than specified purge interval
func (cache *CacheMap) PurgeExpired() {
	ticker := time.NewTicker(cache.PurgeInterval)
	for {
		select {
		case <-ticker.C:
			log.Info("CacheManager - Attempting to delete expired entries")
			for cond, result := range cache.Items {
				if time.Since(result.Timestamp) > cache.PurgeInterval {
					log.Info("CacheManager - Deleting expired entry for ", cond)
					cache.Locker.Lock()
					delete(cache.Items, cond)
					cache.Locker.Unlock()
				}
			}

		}
	}
}

//Set creates an entry in the map if it doesnt exist
// or overwrites the timestamp and retry count if it exists
func (cache *CacheMap) Set(node string, issue string, result ActionResult) {
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
func (cache *CacheMap) GetAll() map[string]ActionResult {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	return cache.Items
}

//GetItem returns value of a given key and whether it exist or not
func (cache *CacheMap) GetItem(key string) (ActionResult, bool) {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	val, found := cache.Items[key]
	if found {
		return val, true
	}
	return ActionResult{}, false
}
