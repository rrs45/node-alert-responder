package controller

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	coreInformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	k8scache "k8s.io/client-go/tools/cache"

	"github.com/box-node-alert-responder/pkg/cache"
)

//AlertWorkerController struct for encapsulating generic Informer methods and pod informer
type AlertWorkerController struct {
	informerFactory   informers.SharedInformerFactory
	podInformer coreInformers.PodInformer
	workerCache *cache.WorkerCache
}

// Run starts shared informers and waits for the shared informer cache to
// synchronize.
func (c *AlertWorkerController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so far
	c.informerFactory.Start(stopCh)
	
	// wait for the initial synchronization of the local cache.
	if !k8scache.WaitForCacheSync(stopCh, c.podInformer.Informer().HasSynced) {
		return fmt.Errorf("Worker Watcher - Failed to sync informer cache")
	}
	return nil
}

func (c *AlertWorkerController) podAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Infof("Worker Watcher - Received pod add event for %s ", pod.Name)
	c.workerCache.SetNew(pod.Name, pod.Status.PodIP)
	
}

func (c *AlertWorkerController) podUpdate(oldPd, newPd interface{}) {
	newpod := newPd.(*v1.Pod)
	log.Infof("Worker Watcher - Received pod update event for %s ", newpod.Name)
	//c.workerCache.Set(newpod.Name, newpod.Status.PodIP)
}

func (c *AlertWorkerController) podDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Infof("Worker Watcher - Received pod delete event for %s ", pod.Name)
    c.workerCache.DelItem(pod.Name)	
}

//NewAlertWorkerController creates a initializes AlertWorkerController struct
//and adds event handler functions
func NewAlertWorkerController(informerFactory informers.SharedInformerFactory, workerCache *cache.WorkerCache) *AlertWorkerController {
	podInf := informerFactory.Core().V1().Pods()

	c := &AlertWorkerController{
		informerFactory:   informerFactory,
		podInformer: podInf,
		workerCache: workerCache,
	}
	podInf.Informer().AddEventHandler(
		// Your custom resource event handlers.
		k8scache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.podAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.podUpdate,
			// Called on resource deletion.
			DeleteFunc: c.podDelete,
		},
	)
	return c
}

//WorkerWatcherStart starts the controller
func WorkerWatcherStart(clientset *kubernetes.Clientset, workerNamespace string, workerCache *cache.WorkerCache) {

	//Set logrus
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.Info("Worker Watcher - Creating informer factory for node-alert-responder to watch in ", workerNamespace, " WorkerNamespace.")
	//Create shared cache informer which resync's every 24hrs
	factory := informers.NewFilteredSharedInformerFactory(clientset, time.Hour*24, workerNamespace,	func(opt *metav1.ListOptions) {} )
	controller := NewAlertWorkerController(factory, workerCache)
	stop := make(chan struct{})
	defer close(stop)
	err := controller.Run(stop)
	if err != nil {
		log.Error("Worker Watcher - Could not run controller :", err)
	}
	select {}
}
