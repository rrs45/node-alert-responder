package controller

import (
	"fmt"
	_ "os"
	_ "os/exec"
	"time"

	"github.com/box-node-alert-responder/pkg/cache"
	"github.com/box-node-alert-responder/pkg/types"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

//Remediate kicks off remediation
func Remediate(client *kubernetes.Clientset, resultsCache *cache.ResultsCache, progressCache *InProgressCache, alertCh <-chan []types.AlertAction) {
	frequency, err := time.ParseDuration("30s")
	if err != nil {
		log.Fatal("Updater - Could not parse interval: ", err)
	}
	ticker := time.NewTicker(frequency)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Responder - \n", resultsCache.GetAll())
		default:
			select {
			case r := <-alertCh:
				//log.Info(r)
				//Send Action result
				for _, item := range r {
					condition := item.Node + "_" + item.Issue
					scheduleFilter(condition, resultsCache, progressCache)
				}
			}
		}
	}
}

func scheduleFilter(condition string, resultsCache *cache.ResultsCache, progressCache *InProgressCache) {
	//Is any worker working on the given node & condition
	if _, working := progressCache.GetItem(condition); !working {
		log.Infof("%s is not currently run by any worker", condition)
		if result, found := resultsCache.GetItem(condition); found {
			log.Infof("%s was worked previously by %s at %v", condition, result.Worker, result.Timestamp)
			if result.Success {
				log.Infof("Last run of %s was successful", condition)
				if result.Timestamp > sucess-wait-interval {
					log.Infof("Last successful run for %s is more than success wait threshold of %v", condition, successWaitPeriod)
					add := progressCache.Set(condition)
					if !add {
						log.Error(condition, " already exists in  Inprogress cache with")
					}
					log.Info("Sending %s condition to scheduler", condition)
					//scheduler(condition)
				} else {
					log.Infof("Last successful run for %s is less than success wait threshold of %v", condition, successWaitPeriod)
					log.Info("Sending %s condition to scheduler", condition)
				}
			} else {
				//If last run failed then check max retries
				log.Infof("Last run of %s failed", condition)
				if result.Retry < max-failed-retry {
					log.Infof("%d failed %d times which is less than %d", condition, result.Retry, maxFailedRetry)
					add := progressCache.Set(condition)
					if !add {
						log.Error(condition, " already exists in  Inprogress cache with")
					}
					log.Info("Sending %s condition to scheduler", condition)
					//scheduler(condition)
				} else {
					log.Infof("%d failed more than %d times hence ignoring", condition, maxFailedRetry)
				}

			}
		} else {
			log.Infof("No record of previous runs for %s in Results cache", condition)
			log.Info("Setting %s in InProgress cache", condition)
			add := progressCache.Set(condition)
			if !add {
				log.Error(condition, " already exists in  Inprogress cache with")
			}
			//No record of
			log.Info("Sending %s condition to scheduler", condition)
			//scheduler(condition)
		}
	}
}
