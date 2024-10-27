package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type MsgPushConfig struct {
	Log struct {
		Path  string
		Level string
	}
	Consul struct {
		Address string
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
}

var Config *MsgPushConfig

func Init() {
	viper.SetConfigName("msgpush")      // name of config file (without extension)
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
