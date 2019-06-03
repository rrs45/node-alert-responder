package options

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"time"
)

type AlertResponderOptions struct {
	ServerAddress  string
	ServerPort     string
	ApiServerHost  string
	LogFile        string
	Namespace      string
	AlertConfigMap string
	PurgeInterval  string
}

//NewAlertResponderOptions returns a flagset
func NewAlertResponderOptions() *AlertResponderOptions {
	return &AlertResponderOptions{}
}

//AddFlags adds options to the flagset
func (aro *AlertResponderOptions) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&aro.ServerAddress, "address", "127.0.0.1",
		"Address to bind the alert generator server.")
	fs.StringVar(&aro.ServerPort, "port", "8090", "Port to bind the alert generator server for /healthz endpoint")
	fs.StringVar(&aro.ApiServerHost, "apiserver-host", "", "Custom hostname used to connect to Kubernetes ApiServer")
	fs.StringVar(&aro.LogFile, "log-file", "/var/log/service/node-alert-generator.log", "Log file to store all logs")
	fs.StringVar(&aro.AlertConfigMap, "alerts-config-map", "npd-alerts", "Name of config map to store alerts")
	fs.StringVar(&aro.Namespace, "namespace", "node-alert-generator", "Namespace where config map will be")
	fs.StringVar(&aro.PurgeInterval, "purge-interval", "48h", "Time period after which cache entries will expire")
}

func (aro *AlertResponderOptions) ValidOrDie() {
	_, err := time.ParseDuration(aro.PurgeInterval)
	if err != nil {
		log.Error("Options - Incorrect alert-update-interval, sample format: 10s or 1m or 1h; ", err)
		log.Panic("Incorrect options")
	}
	/*dir, _ := path.Split(aro.LogFile)
	_, err1 := os.Stat(dir)
	if err1 != nil {
		log.Errorf("Options - Directory does not exist: %s ", dir)
	}
	_, err2 := time.ParseDuration(aro.AlertIgnoreInterval)
	if err2 != nil {
		log.Error("Options - Incorrect alert-ignore-interval, sample format: 10s or 1m or 1h; ", err)
	}
	if err != nil || err1 != nil || err2 != nil {
		log.Panic("Incorrect options")
	} */
}
