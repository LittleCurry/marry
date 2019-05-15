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
	"github.com/LittleCurry/misc/globals"
	"github.com/LittleCurry/misc/driver"
)

func GetUserInfo(c *gin.Context) {

	user_id := c.Query("user_id")

	if len(user_id) == 0 {
		c.JSON(http.StatusOK, err_msg.CodeMsg{1, "user_id不能为空"})
		return
	}
	fmt.Println("user_id:", user_id)
	res := vm.UserRes{"xiaoming", "123123"}
	c.JSON(http.StatusOK, res)

	return

}
func CreateUser(c *gin.Context) {

	createReq := vm.CreateUserReq{}
	err1 := c.Bind(&createReq)
	fmt.Println("creatReq:", createReq)
	if err1 != nil {
		fmt.Println("err1:", err1)
		c.JSON(http.StatusOK, err_msg.CodeMsg{Code: 1, Info: "请求参数格式错误"})
		return
	}

	//if len(createReq.Id) == 0 ||len(createReq.Phone) == 0||len(createReq.Password) == 0{
	//	c.JSON(http.StatusOK,err_msg.CodeMsg{Code:1,Info:"公司名称，或账号，密码不能为空"})
	//	return
	//}
	if len(createReq.NickName) == 0 {
		c.JSON(http.StatusOK, err_msg.CodeMsg{1, "用户名字或者密码不能为空"})
		return
	}
	user := model.User{}
	copier.Copy(&user, createReq)
	user.Passwd = globals.MakeMd5FromString(createReq.Passwd)

	user.UserId = strconv.Itoa(int(time.Now().Unix()))
	user.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	fmt.Println("user.CreateTime:", user.CreateTime)
	_, err2 := driver.MySQL().Insert(&user)
	fmt.Println("err2:", err2)

	c.JSON(http.StatusOK, createReq)
	//onlyUser :=model.User{}
	//driver.MySQL().Where("Id=?",createReq.Id).Where("Phone=?",createReq.Phone).Get(&onlyUser)
	//if len(onlyUser.Id)>0{
	//	c.JSON(http.StatusOK,err_msg.CodeMsg{Code:1,Info:createReq.Phone+"公司账号已经存在"})
	//	return
	//}

}
func UserList(c *gin.Context) {

	users := []model.User{}
	err := driver.MySQL().Find(&users)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, users)

}
func DeleteUser(c *gin.Context) {

}
