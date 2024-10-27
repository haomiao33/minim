package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type OnlineConfig struct {
	Consul struct {
		Address string
	}
	Log struct {
		Path  string
		Level string
	}
	Rpc struct {
		ListenHost string
		ListenPort int
	}
	Kafka struct {
		Addresses     string
		MsgTopic      string
		MsgTopicGroup string
	}
	Redis struct {
		Addr     string
		Password string
	}
	App struct {
		OnlineUserTimeOutSeconds int32
	}
}

var Config *OnlineConfig

func Init() {
	viper.SetConfigName("online")       // name of config file (without extension)
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
