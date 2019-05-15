package config

import (
	"github.com/spf13/viper"
	"fmt"
	"github.com/LittleCurry/misc/helpers"
)

var (
	config GlobalConfig
)

//GlobalConfig 全局配置
type GlobalConfig struct {
	Services []string
	Mails    []string
	Mobiles  []string
}

// Config 返回配置文件
func Config() GlobalConfig {
	return config
}

// ParseConfig 解析配置文件
func ParseConfig(cfg string) {
	viper.SetConfigFile(cfg)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("err17:", err)
		panic(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
}

var AppConf = new(appConf)

type appConf struct {
	helpers.BaseConfig
	RedisAddr string
	MongoAddr string
	DbDsn     string
	BindAddr  string
}

func (this *appConf) GetName() string {
	return "marry"
}
