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

// GpaCrawler is the object representing the crawl engine
// Generally, just use the DefaultGpaCrawler is OK. The difference between different
//GpaCrawlers is just the Configuration.
type GpaCrawler struct {
	config *Configuration
}

func (e *GpaCrawler) SetConfiguration(conf *Configuration) {
	e.config = conf
}

var logger *zap.Logger
var defaultConfig *Configuration

// DefaultGpaCrawler is thread safe. Generally just use this and don't create a new Crawler
var DefaultGpaCrawler *GpaCrawler

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

func init() {
	// 初始化gpa教务系统配置
	defaultConfig = &Configuration{}
	DefaultGpaCrawler = &GpaCrawler{}

	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		logger.Error("Cannot load configuration file automatically: config.json")
		loadDefaultConfiguration()
	} else {
		err = json.Unmarshal(configFile, defaultConfig)
		if err != nil {
			loadDefaultConfiguration()
		}
		// 配置文件错误格式不正确
		if defaultConfig.BaseLocationGPA == nil || defaultConfig.BaseLocationURP == nil {
			logger.Error("config.json 中无法读取所需信息。请正确定义BaseLocationURP:[]string 和 BaseLocationGPA:[]string")
			loadDefaultConfiguration()
		}
	}
	DefaultGpaCrawler.SetConfiguration(defaultConfig)
	logger.Info("Crawler init done.")
}

func loadDefaultConfiguration() {
	logger.Info("Loading default configuration of GpaCrawler", zap.Time("time", time.Now()))
	defaultConfig = &Configuration{
		BaseLocationURP: []string{"http://202.207.247.49", "http://202.207.247.44:8089", "http://202.207.247.51:8065", "http://202.207.247.49"},
		BaseLocationGPA: []string{"http://202.207.247.60/"},
	}
}
