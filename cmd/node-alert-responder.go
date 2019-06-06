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
	cache "github.com/box-node-alert-responder/pkg/cache"
	"github.com/box-node-alert-responder/pkg/controller"
	"github.com/box-node-alert-responder/pkg/controller/types"
)

func initClient(aro *options.AlertResponderOptions) (*kubernetes.Clientset, error) {

	if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".kube", "config")); err == nil {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err)
		}
		return kubernetes.NewForConfig(config)
	} else {
		kubeConfig, err := restclient.InClusterConfig()
		if err != nil {
			panic(err)
		}
		if aro.ApiServerHost != "" {
			kubeConfig.Host = aro.ApiServerHost
		}
		return kubernetes.NewForConfig(kubeConfig)
	}
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
		"NPD-KubeletProxyCertExpiring": "restart_kubeletproxy.yml"}
	log.Info(plays)
	srv := startHTTPServer(aro.ServerAddress, aro.ServerPort)

	var wg sync.WaitGroup
	wg.Add(3)

	// Create an rest client not targeting specific API version
	log.Info("Calling initClient for node-alert-responder")
	clientset, err := initClient(aro)
	if err != nil {
		panic(err)
	}
	log.Info("Successfully generated k8 client for node-alert-responder")
	alertCh := make(chan []types.AlertAction)
	//resultsCh := make(chan map[string]types.ActionResult)
	cache := cache.NewCache(CacheExpireInterval)

	//Watcher
	go func() {
		log.Info("Starting controller for node-alert-responder")
		controller.Start(clientset, aro.Namespace, aro.AlertConfigMap, plays, alertCh)
		log.Info("Watcher - Stopping controller for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()

	//Remediator
	go func() {
		log.Info("Starting remediator for node-alert-responder")
		controller.Remediate(clientset, aro.CacheExpireInterval, aro.ResultsUpdateInterval, cache, alertCh)
		log.Info("Updater - Stopping updater for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()

	//Updater
	go func() {
		log.Info("Starting results configmap updater for node-alert-responder")
		controller.Update(clientset, aro.ResultsNamespace, aro.ResultsConfigMap, aro.ResultsUpdateInterval, cache)
		log.Info("Updater - Stopping updater for node-alert-responder")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not stop http server: %s", err)
		}
		wg.Done()
	}()

	wg.Wait()
}
