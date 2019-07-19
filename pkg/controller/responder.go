package controller

import (
	"time"
	"strconv"
	"strings"

	"github.com/box-node-alert-responder/pkg/types"
	"github.com/box-node-alert-responder/pkg/cache"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

//Remediate kicks off remediation
func Remediate(client *kubernetes.Clientset, resultsCache *cache.ResultsCache, progressCache *cache.InProgressCache, alertCh <-chan types.AlertMap, todo *cache.TodoCache, maxDrained int ) {
	for {
		select {
			case item := <-alertCh:
				log.Debugf("%+v",item)
				retry, err := strconv.Atoi(item.Attr.FailedRetry)
				if err != nil {
					log.Errorf("Responder - Unable convert FailedRetry to int: %v", err)
				}
				run := scheduleFilter(item.NodeCondition, resultsCache, progressCache, item.Attr.SuccessWait, retry, item.Attr.Action, maxDrained )
				if run {	
					log.Infof("Responder - Setting %s condition in Todo cache", item.NodeCondition)
					log.Infof("Responder - progress cache count: %v", progressCache.Count())
					todo.Set(item.NodeCondition, types.Todo{
						Timestamp: time.Now(),
						Action: item.Attr.Action,
						Params: item.Attr.Params,
						})
				}
		} 
	}
}

func scheduleFilter(condition string, resultsCache *cache.ResultsCache, progressCache *cache.InProgressCache, successWaitInterval string, maxRetry int, action string, maxDrained int) bool {
	if resultsCache.GetFailedNodeCount() > maxDrained {
		log.Infof("Responder - Maximum number of drained nodes (%d) reached, ignoring %s", maxDrained, condition)
		return false
	}
	waitDur, _ := time.ParseDuration(successWaitInterval)
	//Is any worker working on the given node & condition
	alert := strings.Split(condition, "_")
	if _, nodePresent := progressCache.GetNode(alert[0]); !nodePresent{
		if _, working := progressCache.GetCondition(alert[0], alert[1]); !working { 
			log.Infof("Responder - %s is not currently run by any worker", condition)
			if result, found := resultsCache.GetItem(condition); found {
				log.Infof("Responder - %s was worked previously by %s at %v", condition, result.Worker, 	result.Timestamp)
				if result.Success {
					log.Infof("Responder - Last run of %s was successful", condition)
					if time.Since(result.Timestamp) > waitDur {
						log.Infof("Responder - Last successful run for %s is more than success wait threshold of %v", 	condition, waitDur)
						return true
					} 
				log.Infof("Responder - Last successful run for %s is less than success wait threshold of %v, ignoring", condition,	 waitDur)
				return false
				} 
				//If last run failed then check max retries
				log.Infof("Responder - Last run of %s failed", condition)
				if result.Retry < maxRetry {
					log.Infof("%s failed %d times which is less than maxRetry:%d ", condition, result.Retry, maxRetry)
					return true
				} else if result.ActionName != action {
					log.Infof("New action defined: %s , old action: %s , for condition: %s", action, result.ActionName, 	condition)
					return true
				}
			
				log.Infof("Responder - %s failed more than %d times and no new action defined hence ignoring", 	condition, maxRetry)
				return false

			} 
			log.Infof("Responder - No record of previous runs for %s in Results cache", condition)
			return true	
		}
		log.Infof("Responder - %s is already being worked on hence ignoring", condition)
		return false
	}
	log.Infof("Responder - %s is already being worked on hence ignoring", condition)
	return false

}
