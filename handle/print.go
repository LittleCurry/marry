package handle

import (
	"github.com/jinzhu/copier"
	"github.com/gin-gonic/gin"
	"fmt"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/err_msg"
	"net/http"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/model"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/driver"
	"gitlab.com/SiivaVideoStudio/cloud_server/me_photo/vm"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/globals"
	"encoding/json"
)

func PrintList(c *gin.Context) {

	activity_id := c.Query("activity_id")
	fmt.Println("activity_id:", activity_id)
	if len(activity_id) == 0 {
		c.JSON(http.StatusOK, err_msg.CodeMsg{1, "activity_id不能为空"})
		return
	}

	printList := []model.Print{}
	driver.MySQL().Where("`activity_id` = ?", activity_id).Find(&printList)
	printListRes := make([]vm.PrintRes, 0)
	copier.Copy(&printListRes, printList)

	c.JSON(http.StatusOK, printListRes)
	return

}

func PrintImg(c *gin.Context) {

	taskId := c.Query("taskId")
	print_id := c.Query("print_id")
	activity_id := c.Query("activity_id")

	if len(taskId) <= 0 || len(print_id) <= 0 || len(activity_id) <= 0 {
		c.JSON(http.StatusBadRequest, err_msg.CodeMsg{1, "taskId, print_id, activity_id都不能为空"})
		return
	}
	conn := globals.GetSocketMap(print_id)
	if conn != nil {
		sentData := &vm.SentData{"print_img", print_id, "https://siiva-video-public.oss-cn-hangzhou.aliyuncs.com/" + activity_id + "/" + taskId + ".jpg"}
		fmt.Println("sentData:", sentData)
		bytes, _ := json.Marshal(&sentData)
		_, err2 := conn.Write(bytes)
		if err2 != nil {
			fmt.Println("err2:", err2)
		}
		c.JSON(http.StatusOK, "已通知"+print_id+"打印")
		return
	} else {
		fmt.Println(print_id+"未注册")
		c.JSON(http.StatusBadRequest, err_msg.CodeMsg{1, print_id+"未注册"})
		return
	}

}
