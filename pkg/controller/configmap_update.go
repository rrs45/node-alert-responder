package controller

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/box-node-alert-responder/pkg/types"
	"github.com/box-node-alert-responder/pkg/cache"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//Update creates config map if it doesnt exist and
//updates the config map with alerts received from watcher
func Update(client *kubernetes.Clientset, ns string, configMap string, resultsUpdateInterval string, resultsCache *cache.ResultsCache) {
	bufCur := make(map[string]map[string]types.ActionResult)
	buf := make(map[string]string)
	frequency, err := time.ParseDuration(resultsUpdateInterval)
	if err != nil {
		log.Errorf("ConfigMap Updater - Could not parse interval: %v", err)
	}
	ticker := time.NewTicker(frequency)

	//Create config map client
	configmapClient := client.CoreV1().ConfigMaps(ns)
	initConfigMap(configmapClient, configMap, resultsCache)
	for {
		select {
		case <-ticker.C:
			log.Info("ConfigMap Updater - Time to save results cache to configmap: ", configMap)
			bufCur = resultsCache.GetAll()
			//Convert ActionResult type to string
			for cond, result := range bufCur {
				rstr, err := json.Marshal(result)
				if err != nil {
					log.Errorf("ConfigMap Updater - unable to marshal %+v: %e", result, err)
				} else {
					buf[cond] = string(rstr)
				}
			}

			//Create config map
			cm := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: configMap,
				},
				Data: buf,
			}
			log.Info("ConfigMap Updater - Updating configmap with ", len(buf), " entries")
			for count := 0; count < 3; count++ {
				result, err := configmapClient.Update(cm)
				if err != nil {
					if count < 3 {
						log.Debugf("ConfigMap Updater - Could not update configmap tried %d times, retrying after 1000ms: %s", count, err)
						time.Sleep(100 * time.Millisecond)
						continue
					} else {
						log.Errorf("ConfigMap Updater - Could not update configmap after 3 attempts: %s", err)
					}
				} else {
					log.Debug("ConfigMap Updater - Updated configmap ", result)
					break
				}
			}
			buf = make(map[string]string)
			bufCur = make(map[string]map[string]types.ActionResult)
		}
		
	}
}

//Check if config map exists if not create new one
func initConfigMap(configmapClient corev1.ConfigMapInterface, name string, resultsCache *cache.ResultsCache) {
	nroCM, err1 := configmapClient.Get(name, metav1.GetOptions{})
	if err1 != nil {
		log.Infof("ConfigMap Updater - %s configmap not found:%v, creating new one", name, err1)
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}
		for count := 0; count < 3; count++ {
			result, err2 := configmapClient.Create(cm)
			if err2 != nil {
				if count < 3 {
					log.Infof("ConfigMap Updater - Could not create configmap tried %d times, retrying after 1000ms: %v", count, err2)
					time.Sleep(100 * time.Millisecond)
					continue
				} else {
					log.Errorf("ConfigMap Updater - Could not create configmap after 3 attempts: %v", err2)
				}
			} else {
				log.Debug("ConfigMap Updater - Created configmap ", result)
				return
			}
		}
	} else {
		log.Infof("ConfigMap Updater - Found existing configmap:%s , updating results cache", name)
		for node, actions := range nroCM.Data {
			var actionResult map[string]types.ActionResult
			err := json.Unmarshal([]byte(actions), &actionResult )
				if err != nil {
					log.Errorf("ConfigMap Updater - Could not unmarshall into JSON:%v", err)
					return
				}
			for action, result := range actionResult {
				//Populate results cache
				resultsCache.Set(node, action, result)
			}
		}

	}
}
