package handle

import (
	"github.com/gin-gonic/gin"
	"encoding/xml"
	"fmt"
	"net/http"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/driver"
	"gopkg.in/mgo.v2/bson"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/model"
	"gitlab.com/SiivaVideoStudio/cloud_server/me_photo/vm"
	"time"
	"strconv"
	"io/ioutil"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/globals"
	"sort"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"github.com/jinzhu/copier"
)

const (
	appid       = "wxb4156cd7e9c16be9"
	mch_id      = "1502192181"
	private_key = "D0A3F589989F8C38226A177A2BF6BB82"
)

func WxchatPayNotify(c *gin.Context) {

	notifyRes := vm.NotifyRes{}
	xml.NewDecoder(c.Request.Body).Decode(&notifyRes)

	reponse := vm.ReturnNotifyRes{}
	orderId := notifyRes.OutTradeNo

	order := model.Order{}
	err1 := driver.Mongo().C("orders").Find(bson.M{"order_id": orderId}).One(&order)
	if err1 != nil {
		fmt.Println("err1:", err1)
	}
	if len(order.OrderId) == 0 {
		reponse.ReturnCode = "FAIL"
		reponse.ReturnMsg = "订单不存在错误"
		fmt.Println("订单不存在错误")
		c.XML(http.StatusOK, reponse)
		return
	}

	err2 := driver.Mongo().C("orders").Update(bson.M{"order_id": orderId}, bson.M{"$set": bson.M{"is_pay": 1}})
	if err2 != nil {
		fmt.Println("err2:", err2)
		reponse.ReturnCode = "FAIL"
		reponse.ReturnMsg = "后台错误"
		c.XML(http.StatusOK, reponse)
		return
	}
	reponse.ReturnCode = "SUCCESS"
	reponse.ReturnMsg = "OK"

	task := model.Task{}
	driver.Mongo().C("tasks").Find(bson.M{"task.taskId": order.Taskid}).Limit(1).One(&task)


	/* 更新wx_scan_bind_task的is_pay为1 */
	if len(order.Taskid) > 0 {

		// err = db.MySQL().In("status=1").In("locker_id", lockerIds).OrderBy("profile_id desc").Find(&lockerParkings)
		//car_update_sql := fmt.Sprintf("update car_user set money=money+%d where id=%d limit 1;", parking.Money, parking.CarUserId)
		//uur, err := session.Exec(car_update_sql)
		// select * from locker where short_id in (56082, 52912);
		//driver.MySQL().Update().In()

		//sql := "update `wx_scan_bind_task` set is_pay=1 where taskId=? limit 1"
		if task.Task != nil && task.Task.(bson.M)["mode"] != nil && task.Task.(bson.M)["mode"] != "photo" {
			sql := "update `wx_scan_bind_task` set is_pay=1 where taskId=? and wx_openid=? limit 1"
			_, err3 := driver.MySQL().Exec(sql, order.Taskid, order.Openid)
			if err3 != nil {
				fmt.Println("err3:", err3)
			}
		}


	}


	if len(task.ActivityId) > 0 {
		activity := model.Activity{}
		driver.Mongo().C("activitys").Find(bson.M{"activity_id": task.ActivityId}).Limit(1).One(&activity)
		oldTaskId := task.Task.(bson.M)["taskId"].(string)

		oldScanList := []model.WxScanBindTask{}
		//driver.MySQL().Where("`wx_openid` = ?", order.Openid).Where("taskId like '" + order.Taskid + "_%'").Find(&oldScanList)
		driver.MySQL().Where("taskId like '" + order.Taskid + "_%'").Find(&oldScanList)

		for i, template := range order.Templates {

			tempScan := model.WxScanBindTask{}
			has, _ := driver.MySQL().Where("activity_id=?", task.ActivityId).Where("wx_openid=?", order.Openid).Where("taskId like '" + order.Taskid + "%'").Where("template_url=?", template).Get(&tempScan)
			if has {
				fmt.Println(order.Taskid+"的"+template+"已经存在了")
				continue
			}

			taskId := oldTaskId + "_" + strconv.Itoa(i+len(oldScanList))
			// 创建多个task
			task.State = "data.ready"
			task.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
			task.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			task.Task.(bson.M)["taskId"] = taskId
			task.Task.(bson.M)["print"] = order.PrintId
			task.Task.(bson.M)["later_setting"] = activity.Settings.(bson.M)["later_setting"]
			task.Task.(bson.M)["later_setting"].(bson.M)["source"].(bson.M)["ground_video"].(bson.M)["src"] = template
			err1 := driver.Mongo().C("tasks").Insert(&task)
			if err1 != nil {
				fmt.Println("err1:", err1)
			}
			// 创建多个wx_scan_bind_task
			wxUser := model.WxScanBindTask{}
			wxUser.WxOpenid = order.Openid
			wxUser.Taskid = taskId
			wxUser.ActivityId = task.ActivityId
			wxUser.TemplateUrl = template
			wxUser.IsPay = 1
			wxUser.CreateTime = time.Now().Format("2006-01-02 15:04:05")
			wxUser.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

			//fmt.Println("wxUser:", wxUser)

			_, err2 := driver.MySQL().Insert(&wxUser)
			if err2 != nil {
				fmt.Println("err2:", err2)
			}
		}
	}

	/* socket通知打印机打印照片*/
	//go func() {
	//	conn := globals.GetSocketMap(order.PrintId)
	//	if conn != nil {
	//		scanList := []model.WxScanBindTask{}
	//		//driver.MySQL().Where("`wx_openid` = ?", order.Openid).Where("`taskId` = ?", order.Taskid).Find(&scanList)
	//		driver.MySQL().Where("`wx_openid` = ?", order.Openid).Where("taskId like '" + order.Taskid + "%'").Find(&scanList)
	//
	//		fmt.Println("len:", len(scanList))
	//
	//		for _, scan := range scanList {
	//
	//			fmt.Println("scan.IsPay:", scan.IsPay)
	//			fmt.Println("order.Print:", order.Print)
	//			fmt.Println("scan.TemplateUrl:", scan.TemplateUrl)
	//			fmt.Println("bool:", globals.Contains(order.Print, scan.TemplateUrl) )
	//
	//			if scan.IsPay == 1 && globals.Contains(order.Print, scan.TemplateUrl) {
	//				sentData := &vm.SentData{"print_img", order.PrintId, "https://siiva-video-public.oss-cn-hangzhou.aliyuncs.com/"+order.ActivityId+"/"+scan.Taskid+".jpg"}
	//				fmt.Println("sentData:", sentData)
	//				bytes, _ := json.Marshal(&sentData)
	//				_, err2 := conn.Write(bytes)
	//				if err2 != nil {
	//					fmt.Println("err2:", err2)
	//				}
	//				time.Sleep(time.Duration(5)*time.Second)
	//			}
	//		}
	//	} else {
	//		fmt.Println(order.PrintId + "未注册")
	//	}
	//}()

	fmt.Println("回复微信")
	c.XML(http.StatusOK, reponse)
	return
}

func generateXml(m map[string]string) string {
	xml := "<xml>"
	for k, v := range m {
		xml += fmt.Sprintf("<%s>%s</%s>", k, v, k)
	}
	xml += "</xml>"
	return xml
}

func callWechatRecharge(orderId string, money string, openId string) (*vm.WechatPayResponse) {

	request := make(map[string]string, 0)

	request["appid"] = appid
	request["body"] = "buy2photos"
	request["mch_id"] = mch_id
	request["nonce_str"] = string(globals.Krand(20, globals.KC_RAND_KIND_LOWER))
	request["notify_url"] = "https://iva.siiva.com/me_photo/pay/notify"
	request["openid"] = openId
	request["out_trade_no"] = orderId
	//request["spbill_create_ip"] = "127.0.0.1" //dev
	request["spbill_create_ip"] = "101.37.151.52" //dev
	request["total_fee"] = money
	request["trade_type"] = "JSAPI"
	request["sign"] = getSign(request, private_key)

	requestStr := generateXml(request)
	fmt.Println("requestStr:", requestStr)

	resp, err1 := http.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", "application/x-www-form-urlencoded", strings.NewReader(requestStr))
	defer resp.Body.Close()
	if err1 != nil {
		fmt.Println("err1:", err1)
	}

	body, err2 := ioutil.ReadAll(resp.Body)
	//fmt.Println("body:", string(body))
	if err2 != nil {
		fmt.Println("err2:", err2)
	}

	response := &vm.WechatPayResponse{}
	err3 := xml.Unmarshal(body, response)
	if err3 != nil {
		fmt.Println("err3:", err3)
	}
	return response
}

func getSign(m map[string]string, key string) string {
	sortedKeys := make([]string, 0)
	signStr := ""
	for k, _ := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		signStr += k + "=" + m[k] + "&"
	}

	signStr += "key" + "=" + key

	h := md5.New()
	fmt.Println(signStr)
	h.Write([]byte(signStr))
	cipherStr := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(cipherStr))
}

func WxchatPay(c *gin.Context) {

	payReq := vm.WxPayReq{}
	err1 := c.Bind(&payReq)
	if err1 != nil {
		fmt.Println("err1:", err1)
	}

	order := model.Order{}
	copier.Copy(&order, payReq)
	order.IsPay = 0
	order.Time = time.Now().Format("2006-01-02 15:04:05")
	order.OrderId = "zhiwei_" + time.Now().Format("20060102150405")
	order.Seller = "me_photo小程序"

	err2 := driver.Mongo().C("orders").Insert(&order)

	if err2 != nil {
		fmt.Println("add_order_err2:", err2)
	}

	wechatPayResponse := callWechatRecharge(order.OrderId, payReq.TotalFee, payReq.Openid)
	response := vm.WxPayRes{}
	response.AppId = wechatPayResponse.Appid
	response.TimeStamp = strconv.Itoa(int(time.Now().Unix()))
	response.NonceStr = wechatPayResponse.NonceStr
	response.Package = "prepay_id=" + wechatPayResponse.PrepayId
	response.SignType = "MD5"

	/* 把微信返回的数据再加密获取一次签名 */
	sm := make(map[string]string, 0)
	sm["appId"] = response.AppId
	sm["timeStamp"] = response.TimeStamp
	sm["nonceStr"] = response.NonceStr
	sm["package"] = response.Package
	sm["signType"] = response.SignType
	response.Sign = getSign(sm, private_key)

	c.JSON(http.StatusOK, response)
	return
}
