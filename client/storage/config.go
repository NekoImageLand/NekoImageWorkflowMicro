package storage

import (
	"NekoImageWorkflowMicro/common"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var configPath string
var configFileName = "NekoImageWorkflowClientConfig"
var configFileNameWithExtension = "NekoImageWorkflowClientConfig.json"

func LoadConfig(info *common.ClientConfig) {
	var config common.ConfigWrapper
	exe, err := os.Executable()
	configPath = filepath.Dir(exe)
	if err != nil {
		logrus.Error("Error getting current directory: %s\n", err)
		return
	}
	if _, err := os.Stat(filepath.Join(configPath, configFileNameWithExtension)); os.IsNotExist(err) {
		CreateConfig()
	} else {
		viper.SetConfigName(configFileName)
		viper.AddConfigPath(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			logrus.Error("Error reading config file, ", err)
		}
		err = viper.Unmarshal(&config)
		if err != nil {
			logrus.Error("Error unmarshalling config file, ", err)
		}
		*info = config.ClientConfig
	}
}

func CreateConfig() {
	logrus.Warning("Config file not found, creating new one.")
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")
	viper.Set("ClientConfig", common.ClientConfig{
		ClientID:              "example-id",
		ClientName:            "example-name",
		ClientRegisterAddress: "https://example.com/register",
		ConsulAddress:         "https://example-consul.com",
		PostUploadPeriod:      300,
		ScraperList:           []common.ScraperType{common.LocalScraperType, common.APIScraperType},
		ScraperConfig: common.ScraperConfig{
			LocalScraperConfig: common.LocalScraperConfig{
				WatchFolders: []string{"/path/to/watch/folder1", "/path/to/watch/folder2"},
			},
			APIScraperConfig: common.APIScraperConfig{
				APIScraperSource: []common.APIScraperSourceConfig{
					{
						APIAddress:           "https://example.com/api",
						ParserJavaScriptFile: "example-parser.js",
					},
				},
			},
		},
	})
	err := viper.SafeWriteConfig()
	if err != nil {
		var configFileAlreadyExistsError viper.ConfigFileAlreadyExistsError
		if errors.As(err, &configFileAlreadyExistsError) {
			logrus.Error("In CreateConfig(), Config file already exists.")
		}
	}
	logrus.Warning("Restart the program to load the new config.")
	os.Exit(114514)
}
