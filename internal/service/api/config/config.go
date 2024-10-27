package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ApiConfig struct {
	Log struct {
		Path  string
		Level string
	}
	Consul struct {
		Address string
	}
	Redis struct {
		Addr     string
		Password string
	}
	Kafka struct {
		Addresses         string
		MsgPartitionCount int
		MsgTopic          string
		MsgTopicGroup     string
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
	App struct {
		Listener string
	}
}

var Config *ApiConfig

func Init() {
	viper.SetConfigName("api")          // name of config file (without extension)
	viper.SetConfigType("toml")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("../config")    // call multiple times to add many search paths
	viper.AddConfigPath("../../config") // call multiple times to add many search paths
	viper.AddConfigPath("./config")     // call multiple times to add many search paths
	viper.AddConfigPath(".")            // optionally look for config in the working directory
	err := viper.ReadInConfig()         // Find and read the config file
	if err != nil {                     // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
}
