package cache

import (
	"sync"
	"time"

	"github.com/box-node-alert-responder/pkg/types"
	log "github.com/sirupsen/logrus"
)

//InProgressCache is a struct to store InProgress of remediation
type InProgressCache struct {
	Items               map[string]map[string]types.InProgress 
	CacheExpireInterval time.Duration
	Locker              *sync.RWMutex
}

//NewInProgressCache instantiates and returns a new cache
func NewInProgressCache(cacheExpireInterval string) *InProgressCache {
	interval, _ := time.ParseDuration(cacheExpireInterval)
	return &InProgressCache{
		Items:               make(map[string]map[string]types.InProgress), //{node:{action:{Inprogress}}}
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
			log.Debug("CacheManager - Attempting to delete expired entries")
			cache.Locker.Lock()
			for node, actionItem := range cache.Items {
				for action, params := range actionItem {
					if time.Since(params.Timestamp) > cache.CacheExpireInterval {
						log.Debug("CacheManager - Deleting expired entry for ", action)
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

//Set adds entry to the cache
func (cache *InProgressCache) Set(node string, condition string, action types.InProgress) {
	cache.Locker.Lock()
	defer cache.Locker.Unlock()
	cond, ok := cache.Items[node]
	if !ok {
		cond = map[string]types.InProgress{condition: action}
		cache.Items[node] = cond
	}
	cache.Items[node][condition] = action
	log.Infof("Progress cache - Setting %+v", cache.Items)
}

//GetAll returns current entries of a cache
func (cache *InProgressCache) GetAll() map[string]map[string]types.InProgress {
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

//GetAction returns value of a given key and whether it exist or not
func (cache *InProgressCache) GetAction(node string, action string) (types.InProgress, bool) {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	_, found := cache.Items[node]
	if found {
		actionVal, actionFound := cache.Items[node][action]
		if 	actionFound {
			return actionVal, true
		}
	}
	return types.InProgress{}, false
}

//GetNode returns value of a given node and whether it exist or not
func (cache *InProgressCache) GetNode(node string) (map[string]types.InProgress, bool) {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	val, found := cache.Items[node]
	if found {
		return val, true
	}
	return nil, false
}

//DelItem deletes a cache item with a given key
func (cache *InProgressCache) DelItem(node string, action string)  {
	cache.Locker.Lock()
	delete(cache.Items[node],action)
	if len(cache.Items[node]) == 0 {
		delete(cache.Items, node)
	}
	log.Infof("Progress cache - deleting %v", cache.Items)
	cache.Locker.Unlock()
}


