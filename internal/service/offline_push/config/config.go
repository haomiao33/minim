package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// [Xiaomi]
// AppSecret = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
//
// [Huawei]
// OAuthClientId = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
// OAuthClientSecret = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
//
// [Vivo]
// AppId = "xxxxxxxxxxxxxxxxxxxxxxxx"
// AppKey = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
// AppSecret = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
//
// [Oppo]
// AppKey = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
// AppServerSecret = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
type OffLineConfig struct {
	XiaoMi struct {
		AppSecret string
	}
	HuaWei struct {
		OAuthClientId     string
		OAuthClientSecret string
	}
	Vivo struct {
		AppId     string
		AppKey    string
		AppSecret string
	}
	Oppo struct {
		AppKey          string
		AppServerSecret string
	}
	Consul struct {
		Address string
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
	Log struct {
		Path  string
		Level string
	}
	Rpc struct {
		ListenHost string
		ListenPort int
	}
}

var Config *OffLineConfig

func Init() {
	viper.SetConfigName("offline")      // name of config file (without extension)
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
