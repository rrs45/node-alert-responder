package controller

import (
	"time"

	"github.com/box-node-alert-responder/pkg/types"
	"github.com/box-node-alert-responder/pkg/cache"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

//Remediate kicks off remediation
func Remediate(client *kubernetes.Clientset, resultsCache *cache.ResultsCache, progressCache *cache.InProgressCache, alertCh <-chan []types.AlertAction, successWaitInterval string, maxRetry int, todo *cache.TodoCache ) {
	/*frequency, err := time.ParseDuration("30s")
	if err != nil {
		log.Fatal("Updater - Could not parse interval: ", err)
	}
	ticker := time.NewTicker(frequency) */
	for {
		select {
		/*case <-ticker.C:
			fmt.Println("Responder - \n", resultsCache.GetAll())
		default:
			select {*/
			case r := <-alertCh:
				//log.Info(r)
				for _, item := range r {
					condition := item.Node + "_" + item.Condition
					run := scheduleFilter(condition, resultsCache, progressCache, successWaitInterval, maxRetry)
					if run {	
						log.Infof("Responder - Setting %s condition in Todo cache", condition)
						log.Infof("Responder - progress cache count: %v", progressCache.Count())
						todo.Set(condition, types.Todo{
							Timestamp: time.Now(),
							Action: item.Action,
							Params: item.Params,
						})

					}
			    }
		}
	}
}

func scheduleFilter(condition string, resultsCache *cache.ResultsCache, progressCache *cache.InProgressCache, successWaitInterval string, maxRetry int) bool {
	waitDur, _ := time.ParseDuration(successWaitInterval)
	//Is any worker working on the given node & condition
	if _, working := progressCache.GetItem(condition); !working {
		log.Infof("Responder - %s is not currently run by any worker", condition)
		if result, found := resultsCache.GetItem(condition); found {
			log.Infof("Responder - %s was worked previously by %s at %v", condition, result.Worker, result.Timestamp)
			if result.Success {
				log.Infof("Responder - Last run of %s was successful", condition)
				if time.Since(result.Timestamp) > waitDur {
					log.Infof("Responder - Last successful run for %s is more than success wait threshold of %v", condition, waitDur)
					return true
				} 
			log.Infof("Responder - Last successful run for %s is less than success wait threshold of %v", condition, waitDur)
			return false
			} 
			//If last run failed then check max retries
			log.Infof("Responder - Last run of %s failed", condition)
			if result.Retry < maxRetry {
				log.Infof("%s failed %d times which is less than %d", condition, result.Retry, maxRetry)
				return true
				} 
			log.Infof("Responder - %s failed more than %d times hence ignoring", condition, maxRetry)
			return false
			
		} 
		log.Infof("Responder - No record of previous runs for %s in Results cache", condition)
		return true	
	}
	return false
}
