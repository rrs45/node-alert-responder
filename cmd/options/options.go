package options

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"	
)

//NewConfigFromFile parses config file  
func NewConfigFromFile(configFile string) (*viper.Viper, error) {
	//dir, file := filepath.Split(configFile)
	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigFile(configFile)
	//v.AddConfigPath(dir)
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}

//Config defines configuration parameters
type Config struct {
	File string
	KubeAPIURL string
	MaxDrained string
}

//GetConfig returna new config file
func GetConfig() *Config {
	return &Config{}
}

//AddFlags takes config file input
func (c *Config) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.File, "file", "/home/rajsingh/go/src/github.com/box-node-alert-responder/config/config.toml",
		"Configuration file path")
		fs.StringVar(&c.KubeAPIURL, "apiserver-override", "", "URL of the kubernetes api server")
		fs.StringVar(&c.MaxDrained, "max-drained-nodes", "200", "Remediation will pause if cluster has these many drained nodes")
}

//ValidOrDie checks some of the options are valid
func ValidOrDie(aro *viper.Viper) {
	log.Infof("Options - %v",aro.AllSettings())
	_, err := time.ParseDuration(aro.GetString("cache.CacheExpireInterval"))
	if err != nil {
		log.Error("Options - Incorrect cache.CacheExpireInterval, sample format: 10s or 1m or 1h; ", err)
		log.Panic("Incorrect options")
	}
	_, err1 := time.ParseDuration(aro.GetString("results.ResultsUpdateInterval"))
	if err1 != nil {
		log.Error("Options - Incorrect results.ResultsUpdateInterval, sample format: 10s or 1m or 1h; ", err1)
		log.Panic("Incorrect options")
	}

	_, err2 := time.ParseDuration(aro.GetString("general.initial_wait_time"))
	if err2 != nil {
		log.Error("Options - Incorrect general.InitialWaitTime, sample format: 10s or 1m or 1h; ", err1)
		log.Panic("Incorrect options")
	}

}
