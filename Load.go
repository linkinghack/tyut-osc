package tyut_osc

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

// Configuration represents the configuration of the crawlers
type Configuration struct {
	BaseLocationURP []string // The base url of main urp system including "http://"
	BaseLocationGPA []string // The base url of GPA system including "http://"
}

type GpaCrawler struct {
	config *Configuration
}

var logger *zap.Logger
var defaultConfig *Configuration
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
	}
	DefaultGpaCrawler.SetConfiguration(defaultConfig)
	logger.Info("Crawler init done.")
}

func infoWithTime(msg ...string) {
	logger.Info(strings.Join(msg, " "), zap.Time("time", time.Now()))
}
func errorWithTime(msg ...string) {
	logger.Error(strings.Join(msg, " "), zap.Time("time", time.Now()))
}
func warnWithTime(msg ...string) {
	logger.Warn(strings.Join(msg, " "), zap.Time("time", time.Now()))
}
func debugWithTime(msg ...string) {
	logger.Debug(strings.Join(msg, " "), zap.Time("time", time.Now()))
}

func loadDefaultConfiguration() {
	defaultConfig = &Configuration{
		BaseLocationURP: []string{"http://202.207.247.49", "http://202.207.247.44:8089", "http://202.207.247.51:8065", "http://202.207.247.49"},
		BaseLocationGPA: []string{"http://202.207.247.60/"},
	}
}

func (e *GpaCrawler) SetConfiguration(conf *Configuration) {
	e.config = conf
}
