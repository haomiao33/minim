package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type MsgConfig struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
	Redis struct {
		Addr     string
		Password string
	}
	Consul struct {
		Address string
	}
	Rpc struct {
		ListenHost string
		ListenPort int
	}
	Kafka struct {
		Addresses         string
		MsgPartitionCount int
		MsgTopic          string
		MsgTopicGroup     string
	}
}

var Config *MsgConfig

func init() {
	viper.SetConfigName("msg")          // name of config file (without extension)
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
