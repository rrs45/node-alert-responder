package controller

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	coreInformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/box-node-alert-responder/pkg/controller/types"
)

//AlertResponderController struct for encapsulating generic Informer methods and configMap informer
type AlertResponderController struct {
	informerFactory   informers.SharedInformerFactory
	configMapInformer coreInformers.ConfigMapInformer
	plays             map[string]string
	alertch           chan<- []types.AlertAction
}

// Run starts shared informers and waits for the shared informer cache to
// synchronize.
func (c *AlertResponderController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so far
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.configMapInformer.Informer().HasSynced) {
		return fmt.Errorf("Failed to sync informer cache")
	}
	return nil
}

func (c *AlertResponderController) configMapAdd(obj interface{}) {
	configMap := obj.(*v1.ConfigMap)
	log.Infof("Watcher - Received configMap add event for %s in watcher.go ", configMap.Name)
}

func playFilter(cm map[string]string, plays map[string]string) []types.AlertAction {
	var actions []types.AlertAction
	for item, param := range cm {
		alert := strings.Split(item, "_")
		if play, ok := plays[alert[1]]; ok {
			actions = append(actions, types.AlertAction{
				Node:   alert[0],
				Issue:  alert[1],
				Params: param,
				Action: play})
		}
	}
	return actions
}

func (c *AlertResponderController) configMapUpdate(oldCM, newCM interface{}) {
	newconfigMap := newCM.(*v1.ConfigMap)

	c.alertch <- playFilter(newconfigMap.Data, c.plays)
}

func (c *AlertResponderController) configMapDelete(obj interface{}) {
	configMap := obj.(*v1.ConfigMap)
	log.Infof("Watcher - Received configMap delete event for %s in watcher.go", configMap.Namespace)
}

//NewAlertResponderController creates a initializes AlertResponderController struct
//and adds event handler functions
func NewAlertResponderController(informerFactory informers.SharedInformerFactory, plays map[string]string, alertch chan<- []types.AlertAction) *AlertResponderController {
	configMapInf := informerFactory.Core().V1().ConfigMaps()

	c := &AlertResponderController{
		informerFactory:   informerFactory,
		configMapInformer: configMapInf,
		plays:             plays,
		alertch:           alertch,
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

//Start starts the controller
func Start(clientset *kubernetes.Clientset, namespace string, configMap string, plays map[string]string, alertch chan<- []types.AlertAction) {

	//Set logrus
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.Info("Watcher - Creating informer factory for node-alert-responder to watch in ", namespace, " namespace.")
	//Create shared cache informer which resync's every 24hrs
	factory := informers.NewFilteredSharedInformerFactory(clientset, time.Hour*24, namespace,
		func(opt *metav1.ListOptions) {
			opt.FieldSelector = fmt.Sprintf("metadata.name=%s", configMap)
		})
	/*factory := informers.NewFilteredSharedInformerFactory(clientset, time.Hour*24, namespace,
	func(opt *metav1.ListOptions) {})*/
	controller := NewAlertResponderController(factory, plays, alertch)
	stop := make(chan struct{})
	defer close(stop)
	err := controller.Run(stop)
	if err != nil {
		log.Error("Watcher - Could not run controller :", err)
	}
	select {}
}
