package main

import (
	"context"
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
	plays := map[string]string{
		"NPD-KubeletIsDown": "raj_test.yml"	}
	
	log.Info(plays)
	srv := startHTTPServer(aro.ServerAddress, aro.ServerPort)

	var wg sync.WaitGroup
	wg.Add(6)

	// Create an rest client not targeting specific API version
	log.Info("Calling initClient for node-alert-responder")
	clientset, err := initClient(aro.APIServerHost)
	if err != nil {
		panic(err)
	}
	log.Info("Successfully generated k8 client for node-alert-responder")
	alertCh := make(chan []types.AlertAction)
	resultsCache := cache.NewResultsCache(aro.CacheExpireInterval)
	inProgressCache := cache.NewInProgressCache(aro.CacheExpireInterval)
	todoCache := cache.NewTodoCache(aro.CacheExpireInterval)
	receiver := controller.NewReceiver(resultsCache, inProgressCache)

	//Watcher
	go func() {
		log.Info("Starting controller for node-alert-responder")
		controller.Start(clientset, aro.AlertsNamespace, aro.AlertConfigMap, plays, alertCh)
		log.Info("Watcher - Stopping controller for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()

	//Remediator
	go func() {
		log.Info("Starting remediator for node-alert-responder")
		controller.Remediate(clientset, resultsCache, inProgressCache , alertCh, aro.WaitAfterSuccess, aro.MaxRetry, todoCache)
		log.Info("Updater - Stopping updater for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()

	//Updater
	go func() {
		log.Info("Starting results configmap updater for node-alert-responder")
		controller.Update(clientset, aro.ResultsNamespace, aro.ResultsConfigMap, aro.ResultsUpdateInterval, resultsCache)
		log.Info("Updater - Stopping updater for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()

	//Scheduler
	go func() {
		log.Info("Starting scheduler for node-alert-responder")
		controller.ScheduleTask(resultsCache, inProgressCache, todoCache, aro.MaxTasks)
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()

	//Receiver
	go func() {
		log.Info("Starting GRPC server for node-alert-responder")
		controller.StartGRPCServer(aro.ReceiverAddress, aro.ReceiverPort, receiver)
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
