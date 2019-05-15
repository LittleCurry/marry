package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gitlab.com/SiivaVideoStudio/cloud_server/me_photo/vm"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/driver"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/err_msg"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func BindOpenid(c *gin.Context) {

	openid := c.Query("openid")
	fmt.Println("openid:", openid)
	if len(openid) == 0 {
		c.JSON(http.StatusOK, err_msg.CodeMsg{1, "openid不能为空"})
		return
	}
	taskId := c.Query("taskId")
	fmt.Println("taskId:", taskId)
	if len(taskId) == 0 {
		c.JSON(http.StatusOK, err_msg.CodeMsg{1, "taskId不能为空"})
		return
	}

	task := model.Task{}
	driver.Mongo().C("tasks").Pipe([]bson.M{
		{"$match": bson.M{"task.taskId": taskId}},
		{"$project": bson.M{"activity_id": 1, "_id": 0}},
	}).One(&task)
	if len(task.ActivityId) == 0 {
		c.JSON(http.StatusOK, err_msg.CodeMsg{1, "未找到该taskId对应的task"})
		return
	}

	bindTask := model.WxScanBindTask{}
	bindTask.WxOpenid = openid
	bindTask.Taskid = taskId
	bindTask.ActivityId = task.ActivityId
	bindTask.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	bindTask.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	tempTask := model.WxScanBindTask{}
	has, _ := driver.MySQL().Where("wx_openid=?", openid).Where("taskId=?", taskId).Get(&tempTask)
	//fmt.Println("has:", has)
	//fmt.Println("bindTask:", bindTask)
	if !has {
		_, err2 := driver.MySQL().Insert(&bindTask)
		if (err2 != nil) {
			fmt.Println("err2:", err2)
			c.JSON(http.StatusOK, err_msg.CodeMsg{1, "绑定失败"})
			return
		}
	}
	c.JSON(http.StatusOK, err_msg.CodeMsg{0, "绑定成功"})
	return

}
func BindOpenidList(c *gin.Context) {
	openid := c.Query("openid")
	//queryTaskid := c.Query("taskId")
	fmt.Println("openid:", openid)
	if len(openid) == 0 {
		c.JSON(http.StatusOK, err_msg.CodeMsg{1, "openid不能为空"})
		return
	}

	openidList := []model.WxScanBindTask{}
	driver.MySQL().Where("`wx_openid` = ?", openid).OrderBy("create_time desc").Find(&openidList)
	//driver.MySQL().Where("`wx_openid` = ?",openid).Find(&openidList).And("taskId=?", queryTaskid).Find(&openid)
	//fmt.Println("openidList", openidList)

	/*转换成客户端显示的格式*/
	openidRes := make([]vm.OpenidRes, 0)
	copier.Copy(&openidRes, openidList)
	for index, temp := range openidRes {

		task := model.Task{}
		driver.Mongo().C("tasks").Find(bson.M{"task.taskId": temp.Taskid}).One(&task)
		//openidRes[index].FileName = task.Task.(bson.M)["fileName1"].(string)
		if task.Task.(bson.M)["files"] != nil {
			files := task.Task.(bson.M)["files"].([]interface{})
			if len(files) > 0 {
				openidRes[index].FileName = files[0].(string)
			}
		}
		if task.Task.(bson.M)["fileName1"] != nil {
			openidRes[index].FileName = task.Task.(bson.M)["fileName1"].(string)
		}
		openidRes[index].State = task.State
	}

	c.JSON(http.StatusOK, openidRes)
	return

}
func TaskInfo(c *gin.Context) {
	taskId := c.Query("taskId")
	query := bson.M{"task.taskId": taskId}
	fmt.Println("taskId:", taskId)
	if len(taskId) == 0 {
		c.JSON(http.StatusOK, err_msg.CodeMsg{1, "taskId不能为空"})
		return
	}
	task := model.Task{}
	driver.Mongo().C("tasks").Find(query).One(&task)
	fmt.Println("task:", task)
	tasksRes := vm.TaskRes{}
	copier.Copy(&tasksRes, task)
	//查找一条记录里的key，并返回给客户端
	if task.Task.(bson.M)["files"] != nil {
		files := task.Task.(bson.M)["files"].([]interface{})
		if len(files) > 0 {
			tasksRes.FileName1 = files[0].(string)
		}
	}
	if task.Task.(bson.M)["fileName1"] != nil {
		tasksRes.FileName1 = task.Task.(bson.M)["fileName1"].(string)
	}
	tasksRes.TaskId = task.Task.(bson.M)["taskId"].(string)
	tasksRes.Mode = "video"
	if task.Task.(bson.M)["mode"] != nil {
		tasksRes.Mode = task.Task.(bson.M)["mode"].(string)
	}

	activity := model.Activity{}
	driver.Mongo().C("activitys").Pipe([]bson.M{
		{"$match": bson.M{"activity_id": task.ActivityId}},
		{"$project": bson.M{"activity_name": 1, "settings": 1, "_id": 0}},
	}).One(&activity)
	if activity.Settings != nil && activity.Settings.(bson.M)["later_setting"] != nil && activity.Settings.(bson.M)["later_setting"].(bson.M)["cut_param"] != nil && activity.Settings.(bson.M)["later_setting"].(bson.M)["cut_param"].(bson.M)["enablegreenscreen"] != nil {
		tasksRes.Enablegreenscreen = activity.Settings.(bson.M)["later_setting"].(bson.M)["cut_param"].(bson.M)["enablegreenscreen"].(bool)
	}

	if tasksRes.Mode == "video" && activity.Settings != nil && activity.Settings.(bson.M)["templates"] != nil {
		tasksRes.OriginalTotal = activity.Settings.(bson.M)["templates"].([]interface{})[0].(bson.M)["original_total"].(int)
		tasksRes.Total = activity.Settings.(bson.M)["templates"].([]interface{})[0].(bson.M)["total"].(int)
		tasksRes.TemplateName = activity.Settings.(bson.M)["templates"].([]interface{})[0].(bson.M)["template_name"].(string)
	}

	c.JSON(http.StatusOK, tasksRes)
	return
}
