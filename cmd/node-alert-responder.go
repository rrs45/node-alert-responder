package main

import (
	"context"
	"time"
	"flag"
	"net/http"
	//"net/http/pprof"
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
		return kubernetes.NewForConfig(kubeConfig)
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
	//mux.HandleFunc("/debug/pprof/", pprof.Index)
    //mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
    //mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
    //mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
    //mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	go func() {
		log.Info("Starting HTTP server for alert-responder")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Could not start http server: %s", err)
		}
	}()
	return srv
}

func main() {
	//Parse command line options
	conf := options.GetConfig()
	conf.AddFlags(flag.CommandLine)
	flag.Parse()
	naro, err := options.NewConfigFromFile(conf.File)
	if err!= nil {
		log.Fatalf("Cannot parse config file: %v", err)
	}
	options.ValidOrDie(naro)
	logFile, _ := os.OpenFile(naro.GetString("general.LogFile"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	defer logFile.Close()
	
	//Set logrus
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	
	srv := startHTTPServer(naro.GetString("server.ServerAddress"), naro.GetString("server.ServerPort"))

	var wg sync.WaitGroup
	wg.Add(7)

	// Create an rest client not targeting specific API version
	log.Info("Calling initClient for node-alert-responder")
	clientset, err := initClient(naro.GetString("kube.APIServerHost"))
	if err != nil {
		panic(err)
	}
	log.Info("Successfully generated k8 client for node-alert-responder")
	alertCh := make(chan types.AlertMap)
	resultsCache := cache.NewResultsCache(naro.GetString("cache.CacheExpireInterval"))
	inProgressCache := cache.NewInProgressCache(naro.GetString("cache.CacheExpireInterval"))
	todoCache := cache.NewTodoCache(naro.GetString("cache.CacheExpireInterval"))
	workerCache := cache.NewWorkerCache()
	receiver := controller.NewReceiver(resultsCache, inProgressCache, workerCache)
	
	//WorkerWatcherStart
	go func() {
		log.Info("Starting worker watcher controller for node-alert-responder")
		controller.WorkerWatcherStart(clientset, naro.GetString("worker.WorkerNamespace"), workerCache, inProgressCache)
		log.Info("Watcher - Stopping worker watcher controller for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()
	log.Infof("Waiting %s for workers to be discovered", naro.GetString("general.InitialWaitTime"))
	time.Sleep(naro.GetDuration("general.InitialWaitTime"))

	//Receiver
	go func() {
		controller.GetWorkerStatus(workerCache, inProgressCache, naro.GetString("worker.WorkerPort"))
		log.Info("Starting GRPC receiver for node-alert-responder")
		controller.StartGRPCServer(naro.GetString("receiver.ReceiverAddress"), naro.GetString("receiver.ReceiverPort"), receiver)
		wg.Done()
	}()

	log.Infof("Waiting %s for workers to be discovered", naro.GetString("general.InitialWaitTime"))
	time.Sleep(naro.GetDuration("general.InitialWaitTime"))

	//Results ConfigMap Updater
	go func() {
		log.Info("Starting results configmap updater for node-alert-responder")
		controller.Update(clientset, naro.GetString("results.ResultsNamespace"), naro.GetString("results.ResultsConfigMap"), naro.GetString("results.ResultsUpdateInterval"), resultsCache)
		log.Info("Updater - Stopping updater for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()
	
	log.Infof("Waiting %s for workers to be discovered", naro.GetString("general.InitialWaitTime"))
	time.Sleep(naro.GetDuration("general.InitialWaitTime"))

	//AlertWatcher
	go func() {
		log.Info("Starting alerts configmap controller for node-alert-responder")
		controller.AlertWatcherStart(clientset, naro.GetString("alerts.AlertsNamespace"), naro.GetString("alerts.AlertConfigMap"), alertCh)
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
		controller.ScheduleTask(workerCache, resultsCache, inProgressCache, todoCache, naro.GetInt("worker.MaxTasks"), naro.GetString("worker.WorkerPort"))
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
