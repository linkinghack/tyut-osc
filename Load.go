package tyut_osc

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Configuration represents the configuration of the crawlers
type Configuration struct {
	BaseLocationURP []string // The base url of main urp system including "http://"
	BaseLocationGPA []string // The base url of GPA system including "http://"
}

var config *Configuration

func init() {
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("Cannot load configuration file automatically: config.json")
		loadDefaultConfiguration()
	} else {
		err = json.Unmarshal(configFile, config)
		if err != nil {
			loadDefaultConfiguration()
		}
	}
}

func loadDefaultConfiguration() {
	config = &Configuration{
		BaseLocationURP: []string{"http://202.207.247.49", "http://202.207.247.44:8089", "http://202.207.247.51:8065", "http://202.207.247.49"},
		BaseLocationGPA: []string{"http://202.207.247.60/"},
	}
}

func SetConfiguration(conf *Configuration) {
	config = conf
}