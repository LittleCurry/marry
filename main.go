package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gitlab.com/SiivaVideoStudio/cloud_server/me_photo/config"
	"gitlab.com/SiivaVideoStudio/cloud_server/me_photo/handle"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/globals"
	"net/http"
	"time"
)

var (
	env        = flag.String("env", "dev", "Running Environment")
	configFile = flag.String("config", "", "The path of configuration file")
)

func init() {
	parseFlag()
	loadConfig()
	initRedis()
	initMongo()
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

	/* 显示屏二维码 */
	router.GET("/qrcode", handle.GetQrcode)
	router.GET("/qrcode_img/:name", handle.QrcodeImg)

	/* user */
	router.GET("/user", handle.GetUserInfo)
	router.POST("/user/creat", handle.CreateUser)
	router.GET("/user/info", handle.GetUserInfo)
	router.GET("user/list", handle.UserList)
	router.POST("user/delete", handle.DeleteUser)

	/*activity*/
	router.GET("/activity", handle.Activity)
	router.GET("/activity/templates", handle.ActivityTemplate)
	router.GET("/activity/info", handle.ActivityInfo)

	/* wechat */

	/* 微信支付回调 */
	router.POST("pay/notify", handle.WxchatPayNotify)
	router.POST("pay", handle.WxchatPay)

	/* task */
	/**
		 * @api {get}  /task/bind_openid绑定taskId
		 * @apiGroup task
		 * @apiVersion 0.1.0
		 * @apiParam {String}  openid 用户id
		 * @apiParam {String}  taskId 照片或视频id
		 *
		 * @apiSuccess {String} code
		 * @apiSuccess {String} info
		 *
		 * @apiSuccessExample {json} Success-Response
		 *  HTTP/1.1 200 OK
		 *  {
    	 *	  "code": 0
    	 *	  "info":"绑定成功"
		 * }
		 *
		 * @apiErrorExample {json} Error-Response
		 *   HTTP/1.1 200 BadRequest
		 *   {
		 *     "code": 1
		 *     "info": "绑定失败"
		 * }
		 */

	/**
			 * @api {get}  /task/bind_openid_list 任务列表BindOpenidList
			 * @apiGroup task
			 * @apiVersion 0.1.0
			 * @apiParam {String}  openid 用户id
			 *
			 * @apiSuccess {String} code
			 * @apiSuccess {object} param.openidList 用户列表
			 * @apiSuccess {String} param.activity_id  活动id
			 * @apiSuccess {String} param.task_id  task_id
			 * @apiSuccess {number} param.is_pay  是否支付(0:未支付；1:已支付)
			 * @apiSuccess {String} param.state  task处理状态
			 * @apiSuccess {String} param.create_time  task创建时间
			 *@apiSuccess {String} param.update_time task更新的时间
			 * @apiSuccessExample {json} Success-Response
			 *  HTTP/1.1 200 OK
			 *  {
	    	 *	  "code": 0
	    	 *	  "result": [ {
			 *			 "id": 1991,
			 *			"activity_id": "",
			 *			"wx_openid": "111",
			 *			"task_id": "",
			 *			"is_pay": 0,
		     *           "create_time": "2019-04-03 16:05:39",
		     *           "update_time": ""
			 *		},
			 *    ...
		     *  ]
			 * }
			 *
			 * @apiErrorExample {json} Error-Response
			 *   HTTP/1.1 200 BadRequest
			 *   {
			 *     "code": 1
			 *     "info": "未找到该openid"
			 * }
	*/
	router.GET("/task/bind_openid", handle.BindOpenid)
	router.GET("/task/bind_openid_list", handle.BindOpenidList)
	router.GET("/task/info", handle.TaskInfo)

	/**
		 * @api {get}  /order/list 订单详情
		 * @apiGroup Order
		 * @apiVersion 0.1.0
		 * @apiParam {String}  openid 用户id
		 *
		 * @apiSuccess {String} code
		 * @apiSuccess {[]Object} result 订单list
		 * @apiSuccess {String} order_id 订单id
		 *
		 * @apiSuccessExample {json} Success-Response
		 *  HTTP/1.1 200 OK
		 *  {
    	 *	  "code": 0
    	 *	  "result": [ {
		 *			"order_id":"123afeeasd",
		 *			"order_state":"complete",
		 *			"total_fee": 300,
		 *			"pay_time":"2018-12-12 12:12:13",
		 *			"goods_info": {"mask":"123123"},
		 *		},
		 *    ...
	     *  ]
		 * }
		 *
		 * @apiErrorExample {json} Error-Response
		 *   HTTP/1.1 200 BadRequest
		 *   {
		 *     "code": 1
		 *     "info": "未找到该openid"
		 * }
		 */

	router.GET("/order/list", handle.OrderList)
	/**
		 * @api {post}  /data/count 数据统计
		 * @apiGroup data
		 * @apiVersion 0.1.0
		 * @apiParam {String}  openid 用户id
		 * @apiParam {String}  taskId 照片或视频id
		 * @apiParam {String}  mode 统计内容（mode：like点赞、download下载、share分享）
		 *
		 * @apiSuccess {String} code
		 * @apiSuccess {String} info
		 *
		 * @apiSuccessExample {json} Success-Response
		 *  HTTP/1.1 200 OK
		 *  {
    	 *	  "code": 0
    	 *	  "info":"统计成功"
		 * }
		 *
		 * @apiErrorExample {json} Error-Response
		 *   HTTP/1.1 200 BadRequest
		 *   {
		 *     "code": 1
		 *     "info": "统计失败"
		 * }
		 */

	/**
		 * @api {get}  /discount/get_discounts 优惠政策
		 * @apiGroup discount
		 * @apiVersion 0.1.0
		 * @apiParam {String}  activity_id 活动id
		 *
		 * @apiSuccess {String} code
         * @apiSuccess {[]Object} result 优惠政策
		 *
		 * @apiSuccessExample {json} Success-Response
		 *  HTTP/1.1 200 OK
		 *  {
    	 *	  "code": 0
    	 *	  "result": ["买三送一"]
		 * }
		 *
		 * @apiErrorExample {json} Error-Response
		 *   HTTP/1.1 200 BadRequest
		 *   {
		 *     "code": 1
		 *     "info": "该活动下未找到优惠政策"
		 * }
		 */

	/**
		 * @api {get}  /templet/get_templates 多模板获取
		 * @apiGroup templet
		 * @apiVersion 0.1.0
		 * @apiParam {String}  activity_id 活动id
		 *
		 * @apiSuccess {String} code
		 * @apiSuccess {[]Object} result 模板list
		 * @apiSuccess {object} param
		 * @apiSuccess {String} param.templet_id 模板id
		 * @apiSuccess {String} param.imgUrl 模板图片路径
		 * @apiSuccess {String} param.type 模板类型(照片或视频)
		 * @apiSuccess {String} param.price 模板价格
		 *
		 * @apiSuccessExample {json} Success-Response
		 *  HTTP/1.1 200 OK
		 *  {
    	 *	  "code": 0
    	 *	  "result": [ {
		 *			"templet_id":"123afeeasd",
		 *			"imgUrl":"https://siiva-video-public.oss-cn-hangzhou.aliyuncs.com/soccer_1541475970gh_1553833274420.jpg",
		 *			"type": 照片,
		 *			"price":"45"
		 *		},
		 *    ...
	     *  ]
		 * }
		 *
		 * @apiErrorExample {json} Error-Response
		 *   HTTP/1.1 200 BadRequest
		 *   {
		 *     "code": 1
		 *     "info": "该活动下未找到模板设定"
		 * }
		 */


	/**
		 * @api {get}  /shooting_point/list 拍摄点列表获取
		 * @apiGroup photo_place
		 * @apiVersion 0.1.0
		 *
		 * @apiSuccess {String} code
		 * @apiSuccess {[]Object} result 拍摄点list
		 * @apiSuccess {object} param
		 * @apiSuccess {String} param.shooting_point  拍摄点名称
		 * @apiSuccess {String} param.shooting_longitude  拍摄点GPS经度
		 * @apiSuccess {String} param.shooting_Latitude  拍摄点GPS纬度
		 * @apiSuccess {String} param.shooting_img  拍摄点横屏缩略图(750px*290px)
		 * @apiSuccess {String} param.shooting_summary  拍摄点介绍说明
		 *
		 * @apiSuccessExample {json} Success-Response
		 *  HTTP/1.1 200 OK
		 *  {
    	 *	  "code": 0
    	 *	  "result": [ {
		 *			"shooting_point":"上海金融中心",
		 *			"shooting_longitude":"121.258302",
		 *			"shooting_Latitude": "30.861681",
		 *			"shooting_img":"https://oss......",
		 *			"shooting_summary": "上海金融中心位于......",
		 *		},
		 *    ...
	     *  ]
		 * }
		 *
		 * @apiErrorExample {json} Error-Response
		 *   HTTP/1.1 200 BadRequest
		 *   {
		 *     "code": 1
		 *     "info": "未找到拍摄点"
		 * }
		 */

	/* 打印机 */
	/**
		 * @api {get}  /print/list 打印机列表获取
		 * @apiGroup print
		 * @apiVersion 0.1.0
		 *
		 * @apiSuccess {object} param
		 * @apiSuccess {String} param.print_id  打印机设备id
		 * @apiSuccess {String} param.address  地址
		 * @apiSuccess {String} param.lon  经度
		 * @apiSuccess {String} param.lat  纬度
		 * @apiSuccess {String} param.mark  描述
		 * @apiSuccess {String} param.create_time  创建时间
		 *
		 * @apiSuccessExample {json} Success-Response
		 *  HTTP/1.1 200 OK
		 *  [ {
		 *			"print_id":"AB12ACD3",
		 *			"address":"上海市陆家嘴",
		 *			"lon": "30.861681",
		 *			"lat":"103.234345",
		 *			"mark": "上海金融中心第一台打印机",
		 *			"create_time": "2019-05-09 12:10:10",
		 *		},
		 *    ...
	     *  ]
		 *
		 * @apiErrorExample {json} Error-Response
		 *   HTTP/1.1 200 BadRequest
		 *   {
		 *     "code": 1
		 *     "info": "未找到拍摄点"
		 * }
		 */
	router.GET("print/list", handle.PrintList)
	router.GET("print", handle.PrintImg)

	go handle.StartSocketServer(config.AppConf.BindAddr)

	router.Run(config.AppConf.HttpPort)

}
