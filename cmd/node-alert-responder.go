package main

import (
	"context"
	"time"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"sync"

	//_ "runtime/pprof"

	log "github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/box-node-alert-responder/cmd/options"
	"github.com/box-node-alert-responder/pkg/cache"
	"github.com/box-node-alert-responder/pkg/controller"
	"github.com/box-node-alert-responder/pkg/types"
)

func initClient(apiserver string) (*kubernetes.Clientset, error) {

	if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".kube", "config")); err == nil {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err)
		}
		return kubernetes.NewForConfig(config)
	} 
	kubeConfig, err := restclient.InClusterConfig()
	if err != nil {
		panic(err)
	}
	if apiserver != "" {
		kubeConfig.Host = apiserver
	}
	return kubernetes.NewForConfig(kubeConfig)
	
}

func startHTTPServer(addr string, port string) *http.Server {
	mux := http.NewServeMux()
	srv := &http.Server{Addr: addr + ":" + port, Handler: mux}
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	go func() {
		log.Info("Starting HTTP server for alert-responder")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Could not start http server: %s", err)
		}
	}()
	return srv
}

func main() {
	//Set logrus
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	//Parse command line options
	aro := options.NewAlertResponderOptions()
	aro.AddFlags(flag.CommandLine)
	flag.Parse()
	aro.ValidOrDie()
	/*f, _ := os.OpenFile(aro.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	log.SetOutput(f) */
	plays := map[string]types.ActionMap{
		"NPD-PuppetRunDelayed": types.ActionMap{
		 Action: "test1.yml",
		SuccessWait: "1m",
		FailedRetry: 2,	},
		"NPD-ChronyIssue": types.ActionMap{
			Action: "test1.yml",
		   SuccessWait: "1m",
		   FailedRetry: 2,	} }
	
	log.Info(plays)
	srv := startHTTPServer(aro.ServerAddress, aro.ServerPort)

	var wg sync.WaitGroup
	wg.Add(7)

	// Create an rest client not targeting specific API version
	log.Info("Calling initClient for node-alert-responder")
	clientset, err := initClient(aro.APIServerHost)
	if err != nil {
		panic(err)
	}
	log.Info("Successfully generated k8 client for node-alert-responder")
	alertCh := make(chan []types.AlertMap)
	resultsCache := cache.NewResultsCache(aro.CacheExpireInterval)
	inProgressCache := cache.NewInProgressCache(aro.CacheExpireInterval)
	todoCache := cache.NewTodoCache(aro.CacheExpireInterval)
	workerCache := cache.NewWorkerCache()
	receiver := controller.NewReceiver(resultsCache, inProgressCache, workerCache)
	
	//WorkerWatcherStart
	go func() {
		log.Info("Starting worker watcher controller for node-alert-responder")
		controller.WorkerWatcherStart(clientset, aro.WorkerNamespace, workerCache)
		log.Info("Watcher - Stopping worker watcher controller for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()
	log.Info("Sleeping 10 seconds for workers to be discovered")
	time.Sleep(10*time.Second)

	//Receiver
	go func() {
		controller.GetWorkerStatus(workerCache, inProgressCache, aro.WorkerPort)
		log.Info("Starting GRPC receiver for node-alert-responder")
		controller.StartGRPCServer(aro.ReceiverAddress, aro.ReceiverPort, receiver)
		wg.Done()
	}()

	//Results ConfigMap Updater
	go func() {
		log.Info("Starting results configmap updater for node-alert-responder")
		controller.Update(clientset, aro.ResultsNamespace, aro.ResultsConfigMap, aro.ResultsUpdateInterval, resultsCache)
		log.Info("Updater - Stopping updater for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()
	
	log.Info("Sleeping 10 seconds for results cache to be populated")
	time.Sleep(10*time.Second)

	//AlertWatcher
	go func() {
		log.Info("Starting alerts configmap controller for node-alert-responder")
		controller.AlertWatcherStart(clientset, aro.AlertsNamespace, aro.AlertConfigMap, plays, alertCh)
		log.Info("Watcher - Stopping alerts configmap controller for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()

	//Remediation filter
	go func() {
		log.Info("Starting remediation filter for node-alert-responder")
		controller.Remediate(clientset, resultsCache, inProgressCache , alertCh, todoCache)
		log.Info("Updater - Stopping updater for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()

	//Scheduler
	go func() {
		log.Info("Starting scheduler for node-alert-responder")
		controller.ScheduleTask(workerCache, resultsCache, inProgressCache, todoCache, aro.MaxTasks, aro.WorkerPort)
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()

	//Results Cache cleanup
	go func() {
		log.Info("Starting results cache cleaner for node-alert-responder")
		resultsCache.PurgeExpired()
		wg.Done()
	}()


	wg.Wait()
}
