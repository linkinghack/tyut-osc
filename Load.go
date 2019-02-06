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
	TempDir         string
	UrpLoginAttempt int
}

var logger *zap.Logger
var OcrPool *OcrEnginePool

//var defaultConfig *Configuration

// DefaultGpaCrawler is thread safe. Generally just use this and don't create a new Crawler
//var DefaultGpaCrawler *GpaCrawler

func init() {
	// 日志配置初始化
	logconfigRaw, err := ioutil.ReadFile("/tyuter/configs/ZapConfig.json")
	if err != nil {
		log.Fatal("Cannot find ZapConfig.json")
	}
	var cfg zap.Config
	if jsonerr := json.Unmarshal(logconfigRaw, &cfg); jsonerr != nil {
		log.Fatal("ZapConfig.json is illeagle")
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err = cfg.Build()
	if err != nil {
		log.Fatal("zap.Logger build failed.")
	}
	logger.WithOptions(zap.AddCaller())
	defer logger.Sync()
	logger.Info("Zap.Logger build success")

}

// 初始化Ocr pool
func init() {
	OcrPool = NewOcrEnginePool(20, 5)
	engine := OcrPool.Get()
	defer OcrPool.Put(engine)

}

func loadDefaultConfiguration() *Configuration {
	logger.Info("Loading default configuration of GpaCrawler", zap.Time("time", time.Now()))
	defaultConfig := &Configuration{
		BaseLocationURP: []string{"http://202.207.247.44:8089", "http://202.207.247.44:8065", "http://202.207.247.44:8059", "http://202.207.247.44:8064", "http://202.207.247.51:8065", "http://202.207.247.49"},
		BaseLocationGPA: []string{"http://202.207.247.60/"},
		UrpLoginAttempt: 10,
	}
	return defaultConfig
}

func loadConfigFromFile(fileName string) *Configuration {
	defaultConfig := &Configuration{}

	configFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Warn("无法加载GPACrawler配置文件: CrawlerConfig.json")
		defaultConfig = loadDefaultConfiguration()
		logger.Info("已使用默认配置创建GPACrawler")
	} else {
		err = json.Unmarshal(configFile, defaultConfig)
		if err != nil {
			defaultConfig = loadDefaultConfiguration()
		}
		// 配置文件错误格式不正确
		if defaultConfig.BaseLocationGPA == nil || defaultConfig.BaseLocationURP == nil {
			logger.Error("CrawlerConfig.json 中无法读取所需信息。请正确定义BaseLocationURP:[]string 和 BaseLocationGPA:[]string")
			defaultConfig = loadDefaultConfiguration()
			logger.Info("使用默认配置创建GPACrawler", zap.Time("time", time.Now()))
		}
	}

	return defaultConfig
}
