package controller

import (
	"encoding/json"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/box-autoremediation/pkg/controller/types"
	"github.com/box-node-alert-responder/pkg/cache"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//Update creates config map if it doesnt exist and
//updates the config map with alerts received from watcher
func Update(client *kubernetes.Clientset, ns string, configMap string, resultsUpdateInterval string, cache *cache.CacheMap) {
	buf := make(map[string]string)
	frequency, err := time.ParseDuration(resultsUpdateInterval)
	if err != nil {
		log.Fatal("Updater - Could not parse interval: ", err)
	}
	ticker := time.NewTicker(frequency)

	//Create config map client
	configmapClient := client.CoreV1().ConfigMaps(ns)
	initConfigMap(configmapClient, configMap)
	for {
		select {
		case <-ticker.C:
			log.Info("Updater - Time to save results cache to configmap: ", configMap)
			//Convert ActionResult type to string
			for cond, result := range cache.GetAll() {
				result, err := json.Marshal(result)
				if err != nil {
					log.Errorf("Updater - unable to marshal %+v: %e", result, err)
				} else {
					buf[cond] = string(ressult)
				}
			}

			//Create config map
			cm := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: configMap,
				},
				Data: buf,
			}
			log.Info("Updater - Updating configmap with ", len(buf), " entries")
			for count := 0; count < 3; count++ {
				result, err := configmapClient.Update(cm)
				if err != nil {
					if count < 3 {
						log.Infof("Updater - Could not update configmap tried %d times, retrying after 1000ms: %s", count, err)
						time.Sleep(100 * time.Millisecond)
						continue
					} else {
						log.Errorf("Updater - Could not update configmap after 3 attempts: %s", err)
					}
				} else {
					log.Debug("Updater - Updated configmap ", result)
					break
				}
			}
			buf = make(map[string]string)
		}
	}
}

func initConfigMap(configmapClient corev1.ConfigMapInterface, name string) {
	_, err1 := configmapClient.Get(name, metav1.GetOptions{})
	if err1 != nil {
		log.Infof("Updater - %s configmap not found, creating new one", name)
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}
		for count := 0; count < 3; count++ {
			result, err2 := configmapClient.Create(cm)
			if err2 != nil {
				if count < 3 {
					log.Infof("Updater - Could not create configmap tried %d times, retrying after 1000ms: %s", count, err2)
					time.Sleep(100 * time.Millisecond)
					continue
				} else {
					log.Errorf("Updater - Could not create configmap after 3 attempts: %s", err2)
				}
			} else {
				log.Debug("Updater - Created configmap ", result)
				break
			}
		}
	} else {
		log.Infof("Updater - %s configmap already exists", name)
	}
}
