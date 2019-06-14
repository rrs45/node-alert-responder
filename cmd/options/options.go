package options

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"time"
)

//AlertResponderOptions defines options for node-alert-responder
type AlertResponderOptions struct {
	ServerAddress string
	ServerPort    string
	ReceiverAddress string
	ReceiverPort string
	APIServerHost string
	LogFile       string

	AlertsNamespace  string
	ResultsNamespace string

	AlertConfigMap   string
	ResultsConfigMap string

	ResultsUpdateInterval string
	CacheExpireInterval   string

	MaxTasks int
}

//NewAlertResponderOptions returns a flagset
func NewAlertResponderOptions() *AlertResponderOptions {
	return &AlertResponderOptions{}
}

//AddFlags adds options to the flagset
func (aro *AlertResponderOptions) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&aro.ServerAddress, "address", "127.0.0.1", "Address to bind the alert generator server.")
	fs.StringVar(&aro.ServerPort, "port", "8090", "Port to bind the alert generator server for /healthz endpoint")
	fs.StringVar(&aro.ReceiverAddress, "receiver-address", "127.0.0.1", "Address to bind the alert generator server.")
	fs.StringVar(&aro.ReceiverPort, "receiver-port", "50040", "Port to bind the alert generator server for /healthz endpoint")
	fs.StringVar(&aro.APIServerHost, "apiserver-host", "", "Custom hostname used to connect to Kubernetes ApiServer")
	fs.StringVar(&aro.LogFile, "log-file", "/var/log/service/node-alert-responder.log", "Log file to store all logs")

	fs.StringVar(&aro.AlertsNamespace, "alert-namespace", "node-alert-generator", "Namespace where alerts config map will be watched from")
	fs.StringVar(&aro.ResultsNamespace, "result-namespace", "node-alert-responder", "Namespace where results config map will be created")

	fs.StringVar(&aro.AlertConfigMap, "alerts-config-map", "npd-alerts", "Name of config map to store alerts")
	fs.StringVar(&aro.ResultsConfigMap, "results-config-map", "nro-results", "Name of config map to store results from remediation")

	fs.StringVar(&aro.ResultsUpdateInterval, "results-update-interval", "30s", "Time period after which cache entries will expire")
	fs.StringVar(&aro.CacheExpireInterval, "cache-expire-interval", "48h", "Time period after which cache entries will expire")

	fs.IntVar(&aro.MaxTasks, "max-tasks", 2, "Maximum number of concurrent remediation tasks at any time")

}

//ValidOrDie checks some of the options are valid
func (aro *AlertResponderOptions) ValidOrDie() {
	_, err := time.ParseDuration(aro.CacheExpireInterval)
	if err != nil {
		log.Error("Options - Incorrect cache-expire-interval, sample format: 10s or 1m or 1h; ", err)
		log.Panic("Incorrect options")
	}
	_, err1 := time.ParseDuration(aro.ResultsUpdateInterval)
	if err1 != nil {
		log.Error("Options - Incorrect results-update-interval, sample format: 10s or 1m or 1h; ", err1)
		log.Panic("Incorrect options")
	}

}
