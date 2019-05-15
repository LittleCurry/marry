package handle

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/err_msg"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/model"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/driver"
	"gopkg.in/mgo.v2/bson"
	"gitlab.com/SiivaVideoStudio/cloud_server/me_photo/vm"
)

func Activity(c *gin.Context) {

}

func ActivityTemplate(c *gin.Context) {

	activity_id := c.Query("activity_id")
	openid := c.Query("openid")
	taskId := c.Query("taskId")
	fmt.Println("activity_id:", activity_id)
	fmt.Println("openid:", openid)
	fmt.Println("taskId:", taskId)
	if len(activity_id) <= 0 {
		c.JSON(http.StatusBadRequest, err_msg.CodeMsg{1, "activity_id不能为空"})
		return
	}

	activity := model.Activity{}
	driver.Mongo().C("activitys").Find(bson.M{"activity_id": activity_id}).One(&activity)
	res := []vm.TemplateRes{}

	if activity.Settings != nil && activity.Settings.(bson.M)["templates"] != nil {

		templates := activity.Settings.(bson.M)["templates"].([]interface{})
		for _, template := range templates {
			temp := vm.TemplateRes{}
			temp.OriginalTotal = template.(bson.M)["original_total"].(int)
			temp.Total = template.(bson.M)["total"].(int)
			temp.TemplateName = template.(bson.M)["template_name"].(string)
			temp.Url = template.(bson.M)["url"].(string)

			tempTask := model.WxScanBindTask{}
			has, _ := driver.MySQL().Where("wx_openid=?", openid).Where("taskId like '" + taskId + "%'").Where("template_url=?", template.(bson.M)["url"].(string)).Get(&tempTask)
			if has {
				temp.IsPay = true
			}
			res = append(res, temp)
		}
	}

	c.JSON(http.StatusOK, res)
	return

}


func ActivityInfo(c *gin.Context) {

	activity_id := c.Query("activity_id")
	activity := model.Activity{}
	driver.Mongo().C("activitys").Find(bson.M{"activity_id": activity_id}).One(&activity)
	c.JSON(http.StatusOK, activity)
	return
}
