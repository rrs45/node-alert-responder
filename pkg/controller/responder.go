package controller

import (
	"time"
	"strconv"

	"github.com/box-node-alert-responder/pkg/types"
	"github.com/box-node-alert-responder/pkg/cache"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

//Remediate kicks off remediation
func Remediate(client *kubernetes.Clientset, resultsCache *cache.ResultsCache, progressCache *cache.InProgressCache, todoCh <-chan types.TodoItem, taskCh chan<- types.TodoItem, maxDrained int ) {
	for {
		select {
			case item := <-todoCh:
				log.Debugf("%+v",item)
				run := scheduleFilter(&item, resultsCache, progressCache, maxDrained )
				if run {	
					log.Debugf("Responder - [node:%s, action:%s] Sending item to task channel", item.Node, item.Action)
					log.Debugf("Responder - progress cache count: %v", progressCache.Count())
					taskCh <- item
				}
		} 
	}
}

func scheduleFilter(item *types.TodoItem, resultsCache *cache.ResultsCache, progressCache *cache.InProgressCache, maxDrained int) bool {
	if resultsCache.GetFailedNodeCount() > maxDrained {
		log.Infof("Responder - [node:%s, action:%s] Ignoring because maximum number of drained nodes (%d) reached", item.Node, item.Action, maxDrained)
		return false
	}
	waitDur, _ := time.ParseDuration(item.SuccessWait)
	//Is any worker working on the given item.Node & condition
	failedRetry, err :=  strconv.Atoi(item.FailedRetry)
	if err != nil {
		log.Fatalf("Responder - [node:%s, action:%s] Could not convert string to int for failed_retry %s", item.Node, item.Action, item.FailedRetry)
	}
	if _, nodePresent := progressCache.GetNode(item.Node); !nodePresent{
		log.Infof("Responder - [node:%s, action:%s] is already being worked on hence ignoring", item.Node, item.Action )
		return false 
	}
	log.Infof("Responder - [node:%s, action:%s] is not currently run by any worker", item.Node, item.Action)
	if result, found := resultsCache.GetItem(item.Node, item.Action);  found {
		log.Infof("Responder - [node:%s, action:%s] was worked previously by %s at %v", item.Node, item.Action, result.Worker, 	result.Timestamp)
		if result.Success {
			log.Infof("Responder - [node:%s, action:%s] Last run was successful", item.Node, item.Action)
			if time.Since(result.Timestamp) > waitDur {
				log.Infof("Responder - [node:%s, action:%s] Last successful run is more than success wait threshold of %v", item.Node, item.Action, waitDur)
				return true
			} 
			log.Infof("Responder - [node:%s, action:%s] Last successful run for %s is less than success wait threshold of %v, ignoring", item.Node, item.Action,	 waitDur)
			return false
		}
		//If last run failed then check max retries
		log.Infof("Responder - [node:%s, action:%s] Last run of failed", item.Node, item.Action)
		if result.Retry <  failedRetry{
			log.Infof("Responder - [node:%s, action:%s] failed %d times which is less than maxFailedRetry:%d ", item.Node, item.Action, result.Retry, failedRetry)
			return true
		}
	
		log.Infof("Responder - [node:%s, action:%s] failed on more than %d times hence ignoring", item.Node, item.Action, failedRetry)
		return false
		}
	log.Infof("Responder - [node:%s, action:%s] No record of previous runs in Results cache", item.Node, item.Action)
	return true	
}

