package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gitlab.com/SiivaVideoStudio/cloud_server/me_photo/vm"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/driver"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/globals"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

func OrderList(c *gin.Context){
	query := bson.M{"del": bson.M{"$ne": "1"}}
	time := c.Query("time")
	is_pay := c.Query("is_pay")
	openid := c.Query("openid")
	activity_id := c.Query("activity_id")
	sort := c.Query("sort")
	start, limit := globals.GetLimitAndStart(c)
	if len(time)>0 {
		query["time"] =bson.M{"$regex": time}
	}
	if len(is_pay)>0{
		intIsPay, _ := strconv.Atoi(is_pay)
		query["is_pay"] = intIsPay
	}
	if len(openid) > 0 {
		query["openid"] = openid
	}
	if len(activity_id) > 0 {
		query["activity_id"] = activity_id
	}

	orders := []model.Order{}
	if len(sort) == 0 {
		sort = "-time"
	}

	err := driver.Mongo().C("orders").Find(query).Sort(sort).Skip(start).Limit(limit).All(&orders)
	if err != nil {
		fmt.Println("err_list:", err)
	}
	res := vm.OrderListRes{}
	count, _ := driver.Mongo().C("orders").Find(query).Count()
	res.Count = count
	copier.Copy(&res.List, orders)

	for _, temp := range res.List {

		task := model.Task{}
		driver.Mongo().C("tasks").Find(bson.M{"task.taskId": temp.Taskid}).One(&task)
		temp.ActivityId = task.ActivityId
		temp.FileName = task.ActivityId+"/"+temp.Taskid + ".mp4"
		if len(task.ActivityId) > 0 && task.Task.(bson.M)["mode"] != nil && task.Task.(bson.M)["mode"].(string) == "photo" {
			temp.FileName = task.ActivityId+"/"+temp.Taskid + ".jpg"
		}


		//temp.ActivityName = task.Task.(map[string]interface{})["fileName1"].(string) // 需要取activity






	}

	/* 求和 */
	sumOrders := []model.Order{}
	query["is_pay"] = 1
	driver.Mongo().C("orders").Find(query).All(&sumOrders)

	total := 0.0
	for _, temp := range sumOrders {
		floatFee, _ := strconv.ParseFloat(temp.TotalFee, 64)
		total += floatFee
	}
	fmt.Println("total:", total)
	res.TotalFees = total

	c.JSON(http.StatusOK, res)
	return

}
