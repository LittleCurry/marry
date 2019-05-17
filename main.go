package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"time"
	"github.com/LittleCurry/misc/globals"
	"github.com/LittleCurry/marry/handle"
	"github.com/LittleCurry/marry/config"
	"github.com/LittleCurry/misc/driver"
	"github.com/LittleCurry/misc/helpers"
	"net/http"
)

var (
	env        = flag.String("env", "dev", "Running Environment")
	configFile = flag.String("config", "", "The path of configuration file")
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

func init() {
	parseFlag()
	loadConfig()
	initOrm()
}

func main() {

	fmt.Println("marry server :", time.Now().Format("2006-01-02 15:04"))
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default() //获得路由实例
	router.Use(globals.Cors())

	// apidoc避免被路由到notfount
	//router.NoRoute(static.Serve("", static.LocalFile("./sign", false)))
	router.Static("/sign", "./sign")
	router.NoRoute(static.Serve("", static.LocalFile("./sign", false)))
	router.NoRoute(static.Serve("/admin", static.LocalFile("./admin", false)))
	version := "version: 0.1.0." + time.Now().Format("20060102.1504")
	router.GET("/version", func(c *gin.Context) { c.String(http.StatusOK, version) })

	/* 微信认证文件 */
	router.GET("/MP_verify_FouzvrsIEmpFd5dU.txt", func(c *gin.Context) { c.File("./resource/MP_verify_FouzvrsIEmpFd5dU.txt"); return })

	/* user */
	router.GET("/user", handle.GetUserInfo)
	router.POST("/user/create", handle.CreateUser)
	router.GET("/user/info", handle.GetUserInfo)
	router.GET("/user/list", handle.UserList)
	router.POST("/user/delete", handle.DeleteUser)

	router.Run(config.AppConf.HttpPort)

}
