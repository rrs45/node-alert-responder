package cache

import (
	"sync"
	"strings"
	"time"
	"container/list"

	"github.com/box-node-alert-responder/pkg/types"

	log "github.com/sirupsen/logrus"
)


//TodoCache is a struct to store results of remediation
type TodoCache struct {
	Items               map[string]types.Todo
	TodoList *list.List
	CacheExpireInterval time.Duration
	Locker *sync.RWMutex
}

//NewTodoCache instantiates and returns a new cache
func NewTodoCache(cacheExpireInterval string) *TodoCache {
	interval, _ := time.ParseDuration(cacheExpireInterval)
	return &TodoCache{
		Items:               make(map[string]types.Todo),
		TodoList: list.New(),
		CacheExpireInterval: interval,
		Locker:              new(sync.RWMutex),
	}
}


//Set creates an entry in the map if it doesnt exist
// or overwrites the timestamp and retry count if it exists
func (cache *TodoCache) Set(cond string, todo types.Todo) {
	cache.Locker.Lock()
	_, found := cache.Items[cond]
	if found {
	cache.Items[cond] = todo
	log.Infof("Todo Cache - Updating %s in todo cache", cond)
	} else {
		log.Infof("Todo Cache - Setting %s in todo cache", cond)
	cache.TodoList.PushBack(cond)
	}
	cache.Locker.Unlock()
}


//GetAll returns current entries of a cache
func (cache *TodoCache) GetAll() map[string]types.Todo {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	return cache.Items
}

//GetItem returns value of a given key and whether it exist or not
func (cache *TodoCache) GetItem() (types.AlertAction, bool) {
	cache.Locker.RLock()
	defer cache.Locker.RUnlock()
	cond := cache.TodoList.Front()
	if cond != nil {
		val := cond.Value.(string)
		n := strings.Split(val, "_")
		return types.AlertAction{
			Node: n[0],
			Condition: n[1],
			Action: cache.Items[val].Action,
			Params: cache.Items[val].Params,
		}, true
	}
	return types.AlertAction{}, false
}

//DelItem deletes an item from Todo list and the map
func (cache *TodoCache) DelItem()  {
	cache.Locker.Lock()
	cond := cache.TodoList.Front()
	log.Infof("Todo Cache - Deleting %s in todo cache", cond.Value.(string))
	cache.TodoList.Remove(cond)
	delete(cache.Items, cond.Value.(string))
	cache.Locker.Unlock()
}