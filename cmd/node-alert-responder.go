package main

import (
	"context"
	"time"
	"flag"
	"net/http"
	"strconv"
	"io"
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
	logFile, _ := os.OpenFile(naro.GetString("general.log_file"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	defer logFile.Close()
	
	//Set logrus
	log.SetFormatter(&log.JSONFormatter{})
	if naro.GetBool("general.debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	
	srv := startHTTPServer(naro.GetString("server.ServerAddress"), naro.GetString("server.ServerPort"))

	var wg sync.WaitGroup
	wg.Add(7)

	// Create an rest client not targeting specific API version
	log.Info("Calling initClient for node-alert-responder")
	clientset, err := initClient(conf.KubeAPIURL)
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
	log.Infof("Waiting %s for workers to be discovered", naro.GetString("general.initial_wait_time"))
	time.Sleep(naro.GetDuration("general.initial_wait_time"))

	//Receiver
	go func() {
		controller.GetWorkerStatus(naro.GetString("certs.cert_file"), naro.GetString("certs.key_file"), naro.GetString("certs.ca_cert_file"), workerCache, inProgressCache, naro.GetString("worker.WorkerPort"), conf.ServerName)
		log.Info("Starting GRPC receiver for node-alert-responder")
		controller.StartGRPCServer(naro.GetString("receiver.ReceiverAddress"), naro.GetString("receiver.ReceiverPort"), naro.GetString("certs.cert_file"), naro.GetString("certs.key_file"), naro.GetString("certs.ca_cert_file"), receiver)
		wg.Done()
	}()

	log.Infof("Waiting %s for workers to get worker status", naro.GetString("general.initial_wait_time"))
	time.Sleep(naro.GetDuration("general.initial_wait_time"))

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
	
	log.Infof("Waiting %s for workers to be discovered", naro.GetString("general.initial_wait_time"))
	time.Sleep(naro.GetDuration("general.initial_wait_time"))

	//AlertWatcher
	go func() {
		log.Info("Starting alerts configmap controller for node-alert-responder")
		controller.AlertWatcherStart(clientset, naro.GetString("alerts.AlertsNamespace"), naro.GetString("alerts.ConfigMapLabel"), alertCh)
		log.Info("Watcher - Stopping alerts configmap controller for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()

	//Remediation filter
	go func() {
		log.Info("Starting remediation filter for node-alert-responder")
		maxDrainedInt, _ := strconv.Atoi(conf.MaxDrained)
		controller.Remediate(clientset, resultsCache, inProgressCache , alertCh, todoCache, maxDrainedInt)
		log.Info("Updater - Stopping updater for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()

	//Scheduler
	go func() {
		log.Info("Starting scheduler for node-alert-responder")
		controller.ScheduleTask(naro.GetString("certs.cert_file"), naro.GetString("certs.key_file"), naro.GetString("certs.ca_cert_file"), workerCache, inProgressCache, todoCache, naro.GetInt("worker.MaxTasks"), naro.GetString("worker.WorkerPort"), conf.ServerName)
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
