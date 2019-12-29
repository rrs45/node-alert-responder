package cache

import (
	"sync"
	log "github.com/sirupsen/logrus"
	"github.com/box-node-alert-responder/pkg/types"
)

//WorkerCache is a struct to store Worker of remediation
type WorkerCache struct {
	Items               map[string]types.Worker	//{podName: {IP, Taskcount, Node}}			
	Locker              *sync.RWMutex
}

//NewWorkerCache instantiates and returns a new cache
func NewWorkerCache() *WorkerCache {
	return &WorkerCache{
		Items:               make(map[string]types.Worker),
		Locker:              new(sync.RWMutex),
	}
}

//SetNew appends entry to the map
func (cache *WorkerCache) SetNew(worker string, ip string, node string) {
	if  _, found := cache.Items[worker]; !found{
		if  ip != "" && node != "" {
			log.Infof("Worker Cache - Setting new worker:%s with IP:%s in cache", worker, ip)
			cache.Locker.Lock()
			cache.Items[worker] = types.Worker{
								IP: ip,
								TaskCount: 0,
								Node: node,
								}
			cache.Locker.Unlock()
		}
	}
}

//Increment appends entry to the map
func (cache *WorkerCache) Increment(worker string) {
    log.Debugf("Worker Cache - Incrementing task count for worker: %s ", worker)
	cache.Locker.Lock()
	v := cache.Items[worker]
	v.TaskCount++
	cache.Items[worker] = v
	cache.Locker.Unlock()
}

//Decrement appends entry to the map
func (cache *WorkerCache) Decrement(worker string) {
    log.Debugf("Worker Cache - Decrementing task count for worker: %s ", worker)
	cache.Locker.Lock()
	v := cache.Items[worker]
	v.TaskCount--
	cache.Items[worker] = v
	cache.Locker.Unlock()
}


//GetNext return an IP of one of the available workers
func (cache *WorkerCache) GetNext(maxTasks int, node string) (string, string, string) {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	for key, val := range cache.Items {
		if val.TaskCount <= maxTasks && val.Node != node{
			return key, val.IP, val.Node
		}
	}
	return "","",""
}

//GetAll returns current entries of a cache
func (cache *WorkerCache) GetAll() map[string]types.Worker {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	return cache.Items
}

//Count returns number of items in cache
func (cache *WorkerCache) Count() int {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	return len(cache.Items)
}

//GetItem returns value of a given key and whether it exist or not
func (cache *WorkerCache) GetItem(key string) (types.Worker, bool) {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	val, found := cache.Items[key]
	if found {
		return val, true
	}
	return val, false
}

//DelItem deletes a cache item with a given key
func (cache *WorkerCache) DelItem(key string)  {
	log.Debugf("Worker Cache - Deleting %s from cache", key)
	cache.Locker.Lock()
	delete(cache.Items,key)
	cache.Locker.Unlock()
}


