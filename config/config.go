package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/sflewis2970/trivia-service/common"
)

const BASE_DIR_NAME string = "trivia-service"
const CONFIG_FILE_NAME = "./config/config.json"
const UPDATE_CONFIG_DATA string = "update"

// Config variable keys
const (
	// System ENV setting
	ENV string = "ENV"

	// Host system info
	HOST string = "HOST"
	PORT string = "PORT"

	// Datastore Server info
	DSNAME string = "DS_NAME"
	DSPORT string = "DS_PORT"

	// Response Messages
	CONGRATS string = "CONGRATS"
	TRYAGAIN string = "TRYAGAIN"
)

// Config variable values
const (
	PRODUCTION string = "PROD"
)

type Message struct {
	CongratsMsg string `json:"congrats"`
	TryAgainMsg string `json:"tryagain"`
}

type ConfigData struct {
	Env           string `json:"env"`
	Host          string `json:"hostname"`
	Port          string `json:"hostport"`
	DataStoreName string `json:"dsname"`
	DataStorePort string `json:"dsport"`
	Messages      Message
}

type config struct {
	cfgData *ConfigData
}

var cfg *config

// Unexported type functions
func (c *config) findBaseDir(currentDir string, targetDir string) int {
	level := 0
	dirs := strings.Split(currentDir, "\\")

	dirsSize := len(dirs)
	for idx := dirsSize - 1; idx >= 0; idx-- {
		if dirs[idx] == targetDir {
			break
		} else {
			level++
		}
	}

	return level
}

func (c *config) readConfigFile() error {
	// Get working directory
	wd, getErr := common.GetWorkingDir()
	if getErr != nil {
		log.Print("Error getting working directory")
		return getErr
	}

	// Find path
	levels := c.findBaseDir(wd, BASE_DIR_NAME)
	for levels > 0 {
		chErr := os.Chdir("..")
		if chErr != nil {
			log.Print("Error changind dir: ", chErr)
		}

		// Update levels
		levels--
	}

	data, readErr := ioutil.ReadFile(CONFIG_FILE_NAME)
	if readErr != nil {
		return readErr
	}

	unmarshalErr := json.Unmarshal(data, c.cfgData)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	return nil
}

func (c *config) getConfigEnv() error {
	// Loading config environment variables
	log.Print("loading config environment variables...")

	// Update config data
	// Base config settings
	c.cfgData.Env = os.Getenv(ENV)
	c.cfgData.Host = os.Getenv(HOST)
	c.cfgData.Port = os.Getenv(PORT)
	c.cfgData.DataStoreName = os.Getenv(DSNAME)
	c.cfgData.DataStorePort = os.Getenv(DSPORT)

	// Get response messages
	c.cfgData.Messages.CongratsMsg = "Congratulations! That is correct"
	c.cfgData.Messages.TryAgainMsg = "Nice Try! Better luck on the next answer"

	return nil
}

func (c *config) GetData(args ...string) (*ConfigData, error) {
	if len(args) > 0 {
		if args[0] == UPDATE_CONFIG_DATA {
			useCfgFile := os.Getenv("USECONFIGFILE")
			if len(useCfgFile) > 0 {
				log.Print("Using config file to load config")

				readErr := cfg.readConfigFile()
				if readErr != nil {
					log.Print("Error reading config file: ", readErr)
					return nil, readErr
				}
			} else {
				log.Print("Using config environment to load config")

				getErr := cfg.getConfigEnv()
				if getErr != nil {
					log.Print("Error getting config environment data: ", getErr)
					return nil, getErr
				}
			}
		}
	}

	return c.cfgData, nil
}

func Get() *config {
	if cfg == nil {
		log.Print("creating config object")

		// Initialize config
		cfg = new(config)

		// Initialize config data
		cfg.cfgData = new(ConfigData)
	}

	log.Print("returning config object")
	return cfg
}
