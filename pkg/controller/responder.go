package controller

import (
	"fmt"
	_ "os"
	_ "os/exec"
	"time"

	"github.com/box-node-alert-responder/pkg/cache"
	"github.com/box-node-alert-responder/pkg/controller/types"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

//Remediate kicks off remediation
func Remediate(client *kubernetes.Clientset, resultsCache *cache.ResultsCache, alertCh <-chan []types.AlertAction) {
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
					resultsCache.Set(item.Node, item.Issue, types.ActionResult{Timestamp: time.Now(),
						ActionName: item.Action,
						Success:    true,
						Retry:      0,
					})
				}
			}
		}
	}
}

/*func runCmd(params string) {
	cmd := exec.Command("ansible-playbook", params)
}*/
