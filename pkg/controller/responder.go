package controller

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	mapcache "github.com/box-node-alert-responder/pkg/cache"
	"github.com/box-node-alert-responder/pkg/controller/types"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

//Remediate kicks off remediation
func Remediate(client *kubernetes.Clientset, purgeInterval string, ch <-chan []types.AlertAction) {
	cache := mapcache.NewCache(purgeInterval)
	frequency, err := time.ParseDuration("30s")
	if err != nil {
		log.Fatal("Updater - Could not parse interval: ", err)
	}
	ticker := time.NewTicker(frequency)
	for {
		select {
		case <-ticker.C:
			fmt.Println(cache.GetAll())
		default:
			select {
			case r := <-ch:
				//log.Info(r)
				//Send Action result
				for _, item := range r {
					cache.Set(item.Node, item.Issue, mapcache.ActionResult{Timestamp: time.Now(),
						ActionName: item.Action,
						Success:    true,
						Retry:      0,
					})
				}
			}
		}
	}
}

func runCmd(params string) {
	cmd := exec.Command("ansible-playbook", params)
}
