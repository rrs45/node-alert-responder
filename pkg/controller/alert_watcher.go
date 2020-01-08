package controller

import (
	"fmt"
	"time"
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	coreInformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/box-node-alert-responder/pkg/types"
)

//AlertResponderController struct for encapsulating generic Informer methods and configMapLabel informer
type AlertResponderController struct {
	informerFactory   informers.SharedInformerFactory
	configMapInformer coreInformers.ConfigMapInformer
	todoCh           chan<- types.TodoItem
}

// Run starts shared informers and waits for the shared informer cache to
// synchronize.
func (c *AlertResponderController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so far
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.configMapInformer.Informer().HasSynced) {
		return fmt.Errorf("Configmap Watcher - Failed to sync informer cache")
	}
	return nil
}

func sendAlert(cm *v1.ConfigMap, ch chan<- types.TodoItem ) {
	for nodeCond, attr := range cm.Data {
		var alertMap types.AlertMap
		err := json.Unmarshal([]byte(attr), &alertMap.Attr )
		if err != nil {
			log.Errorf("Alert Watcher - Could not unmarshall into JSON:%v", err)
		}
		log.Debugf("Alert Watcher - %+v", alertMap.Attr)
		alertMap.NodeCondition = nodeCond
		node, condition := func() (string, string) {
			x := strings.Split(nodeCond, "_")
			return x[0], x[1]
		}()
		todoItem := types.TodoItem{
				Node: node,
				Condition: condition,
				Source: cm.Name,
				Action: alertMap.Attr.Action,
				Params: alertMap.Attr.Params,
				SuccessWait: alertMap.Attr.SuccessWait,
				FailedRetry: alertMap.Attr.FailedRetry,
				Timestamp: alertMap.Attr.Timestamp,
		}
		ch <- todoItem
	}
}

func (c *AlertResponderController) configMapAdd(obj interface{}) {
	configMap := obj.(*v1.ConfigMap)
	log.Infof("ConfigMap Watcher - Received config map add event for %s in watcher.go ", configMap.Name)
	sendAlert(configMap, c.todoCh)
}

func (c *AlertResponderController) configMapUpdate(oldCM, newCM interface{}) {
	newconfigMap := newCM.(*v1.ConfigMap)
	log.Infof("ConfigMap Watcher - Received config map update event for %s in watcher.go ", newconfigMap.Name)
	sendAlert(newconfigMap, c.todoCh)
}

func (c *AlertResponderController) configMapDelete(obj interface{}) {
	configMap := obj.(*v1.ConfigMap)
	log.Infof("ConfigMap Watcher - Received config map delete event for %s in watcher.go", configMap.Name)
}

//NewAlertResponderController creates a initializes AlertResponderController struct
//and adds event handler functions
func NewAlertResponderController(informerFactory informers.SharedInformerFactory, todoCh chan<- types.TodoItem) *AlertResponderController {
	configMapInf := informerFactory.Core().V1().ConfigMaps()

	c := &AlertResponderController{
		informerFactory:   informerFactory,
		configMapInformer: configMapInf,
		todoCh:           todoCh,
	}
	configMapInf.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.configMapAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.configMapUpdate,
			// Called on resource deletion.
			DeleteFunc: c.configMapDelete,
		},
	)
	return c
}

//AlertWatcherStart starts the controller
func AlertWatcherStart(clientset *kubernetes.Clientset, AlertsNamespace string, configMapLabel string, todoCh chan<- types.TodoItem) {

	//Set logrus
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.Info("ConfigMap Watcher - Creating informer factory for node-alert-responder to watch in ", AlertsNamespace, " AlertsNamespace.")
	//Create shared cache informer which resync's every 24hrs
	factory := informers.NewFilteredSharedInformerFactory(clientset, time.Hour*24, AlertsNamespace,
		func(opt *metav1.ListOptions) {
			opt.LabelSelector = fmt.Sprintf("%s", configMapLabel)
		})
	controller := NewAlertResponderController(factory, todoCh)
	stop := make(chan struct{})
	defer close(stop)
	err := controller.Run(stop)
	if err != nil {
		log.Error("ConfigMap Watcher - Could not run controller :", err)
	}
	select {}
}
