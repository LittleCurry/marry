package main

import (
	"flag"
	"fmt"

	"github.com/LittleCurry/marry/config"
	"github.com/LittleCurry/misc/helpers"
	"github.com/LittleCurry/misc/driver"
)

/**
 * 解析参数
 */
func parseFlag() {
	flag.Parse()
	if *configFile == "" {
		fmt.Println("place provide the configuration file path!")
		fmt.Println("Usage:\n --config configpath [--env prod]")
	}
}

/**
 * 载入配置文件
 */
func loadConfig() {

	fmt.Println("config.AppConf:", config.AppConf)
	fmt.Println("env:", *env)
	fmt.Println("configFile:", *configFile)

	helpers.LoadConfigAndSetupEnv(config.AppConf, *env, *configFile)
}

/**
 * 初始化redis
 */
func initRedis() {
	fmt.Println("configs.AppConf.RedisAddr:", config.AppConf.RedisAddr)
	driver.RedisInit(config.AppConf.RedisAddr, 0)
}

/**
 * 初始化redis
 */
func initMongo() {
	fmt.Println("configs.AppConf.MongoAddr:", config.AppConf.MongoAddr)
	driver.MongoInit(config.AppConf.MongoAddr)
}

/**
 * 初始化orm引擎
 */
func initOrm() {
	fmt.Println("configs.AppConf.DbDsn:", config.AppConf.DbDsn)
	driver.OrmInit(config.AppConf.DbDsn)
}

