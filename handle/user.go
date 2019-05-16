package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"strconv"
	"time"
	"github.com/LittleCurry/misc/err_msg"
	"github.com/LittleCurry/marry/vm"
	"github.com/LittleCurry/misc/model"
	"github.com/LittleCurry/misc/driver"
)

func GetUserInfo(c *gin.Context) {

	//user_id := c.Query("user_id")

	//if len(user_id) == 0 {
	//	c.JSON(http.StatusOK, err_msg.CodeMsg{1, "user_id不能为空"})
	//	return
	//}
	//fmt.Println("user_id:", user_id)
	//res := vm.UserRes{"xiaoming", "123123"}
	//c.JSON(http.StatusOK, res)

	return

}
func CreateUser(c *gin.Context) {

	createReq := vm.CreateUserReq{}
	err1 := c.Bind(&createReq)
	fmt.Println("Phone:", createReq.Phone)
	fmt.Println("Address:", createReq.Address)
	fmt.Println("UserName:", createReq.UserName)
	fmt.Println("Introduction:", createReq.Introduction)
	fmt.Println("Birthday:", createReq.Birthday)
	if err1 != nil {
		fmt.Println("err1:", err1)
		c.JSON(http.StatusOK, err_msg.CodeMsg{Code: 1, Info: "请求参数格式错误"})
		return
	}

	user := model.User{}
	copier.Copy(&user, createReq)
	user.UserId = strconv.Itoa(int(time.Now().Unix()))
	user.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	user.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	_, err2 := driver.MySQL().Insert(&user)
	if err2 != nil {
		fmt.Println("err2:", err2)
	}

	c.JSON(http.StatusOK, createReq)
	return

}
func UserList(c *gin.Context) {

	users := []model.User{}
	err := driver.MySQL().Find(&users)
	if err != nil {
		fmt.Println("err:", err)
	}
	usersRes := []vm.UserRes{}
	copier.Copy(&usersRes, users)
	c.JSON(http.StatusOK, usersRes)

}
func DeleteUser(c *gin.Context) {

}
