package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"github.com/LittleCurry/misc/globals"
	"github.com/LittleCurry/marry/handle"
	"github.com/LittleCurry/marry/config"
)

var (
	env        = flag.String("env", "dev", "Running Environment")
	configFile = flag.String("config", "", "The path of configuration file")
)

func init() {
	parseFlag()
	loadConfig()
	initOrm()
}

func main() {

	fmt.Println("me_photo server :", time.Now().Format("2006-01-02 15:04"))
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default() //获得路由实例
	router.Use(globals.Cors())

	// apidoc避免被路由到notfount
	router.NoRoute(static.Serve("/apidoc", static.LocalFile("./apidoc", false)))
	version := "version: 0.1.0." + time.Now().Format("20060102.1504")
	router.GET("/version", func(c *gin.Context) { c.String(http.StatusOK, version) })

	/* 微信认证文件 */
	router.GET("/MP_verify_FouzvrsIEmpFd5dU.txt", func(c *gin.Context) { c.File("./resource/MP_verify_FouzvrsIEmpFd5dU.txt"); return })


	/* user */
	router.GET("/user", handle.GetUserInfo)
	router.POST("/user/creat", handle.CreateUser)
	router.GET("/user/info", handle.GetUserInfo)
	router.GET("user/list", handle.UserList)
	router.POST("user/delete", handle.DeleteUser)

	router.Run(config.AppConf.HttpPort)

}
