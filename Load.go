package tyut_osc

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"log"
	"time"
)

// Configuration represents the configuration of the crawlers
type Configuration struct {
	BaseLocationURP []string // The base url of main urp system including "http://"
	BaseLocationGPA []string // The base url of GPA system including "http://"
}

var logger *zap.Logger

//var defaultConfig *Configuration

// DefaultGpaCrawler is thread safe. Generally just use this and don't create a new Crawler
//var DefaultGpaCrawler *GpaCrawler

func init() {
	// 日志配置初始化
	logconfigRaw, err := ioutil.ReadFile("logConfig.json")
	if err != nil {
		log.Fatal("Cannot find logConfig.json")
	}
	var cfg zap.Config
	if jsonerr := json.Unmarshal(logconfigRaw, &cfg); jsonerr != nil {
		log.Fatal("logConfig.json is illeagle")
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err = cfg.Build()
	if err != nil {
		log.Fatal("zap.Logger build failed.")
	}
	defer logger.Sync()
	logger.Info("Zap.Logger build success")

}

func loadDefaultConfiguration() *Configuration {
	logger.Info("Loading default configuration of GpaCrawler", zap.Time("time", time.Now()))
	defaultConfig := &Configuration{
		BaseLocationURP: []string{"http://202.207.247.49", "http://202.207.247.44:8089", "http://202.207.247.51:8065", "http://202.207.247.49"},
		BaseLocationGPA: []string{"http://202.207.247.60/"},
	}
	return defaultConfig
}
