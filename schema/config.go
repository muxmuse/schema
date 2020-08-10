package schema

import (
	"os"
	"os/user"
	"gopkg.in/yaml.v2"
	"io/ioutil"
  "path/filepath"
  "errors"
  "github.com/muxmuse/schema/mfa"
)

var Config TConfig
var SelectedConnectionConfig TConnectionConfig
var WorkingDirectory string

type TConnectionConfig struct {
	Name string
	Url string
	User string
	Password string
	Selected bool
}

type TConfig struct {
	Connections []TConnectionConfig
}

func getConfig() TConfig {
	currentUser, err := user.Current()
	mfa.CatchFatal(err)
	
	yamlFile, err := ioutil.ReadFile(filepath.Join(
		currentUser.HomeDir, ".schemapm", "config.yaml"))
	mfa.CatchFatal(err)

	var config TConfig
	err = yaml.Unmarshal(yamlFile, &config)
	mfa.CatchFatal(err)

	return config
}

func getSelectedConnectionConfig(config TConfig) TConnectionConfig {
	if len(config.Connections) == 0 {
		errors.New("Please configure at least one connection.")
	}

	for _, cc := range config.Connections {
		if cc.Selected {
			return cc
		}
	}

	return config.Connections[0]
}

func init() {
	Config = getConfig()
	SelectedConnectionConfig = getSelectedConnectionConfig(Config)
	dir, err := os.Getwd()
	mfa.CatchFatal(err)
	WorkingDirectory = dir
}
