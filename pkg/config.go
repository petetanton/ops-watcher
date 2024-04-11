package pkg

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	JiraUsername string   `mapstructure:"jira_username"`
	JiraPassword string   `mapstructure:"jira_password"`
	JiraToken    string   `mapstructure:"jira_token"`
	JiraBaseUrl  string   `mapstructure:"jira_baseurl"`
	JiraEnabled  bool     `mapstructure:"jira_enabled"`
	JiraQuery    []string `mapstructure:"jira_query"`
}

func NewConfig() (*Config, error) {
	viper.SetConfigName("ops-watcher")        // name of config file (without extension)
	viper.SetConfigType("yaml")               // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/.ops-watcher") // call multiple times to add many search paths
	viper.AddConfigPath(".")                  // optionally look for config in the working directory
	err := viper.ReadInConfig()               // Find and read the config file
	if err != nil {                           // Handle errors reading the config file
		return nil, fmt.Errorf("Fatal error config file: %s \n", err)
	}

	config := &Config{}

	err = viper.Unmarshal(&config)

	return config, err
}
