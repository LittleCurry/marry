package handle

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gitlab.com/SiivaVideoStudio/cloud_server/me_photo/vm"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/driver"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/err_msg"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/globals"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/model"
	"net/http"
	"strconv"
	"time"
)

func GetUserInfo(c *gin.Context)  {

	user_id := c.Query("user_id")

	if len(user_id) == 0 {
		c.JSON(http.StatusOK, err_msg.CodeMsg{1,"user_id不能为空"})
		return
	}
	fmt.Println("user_id:", user_id)
	res := vm.UserRes{"xiaoming", "123123"}
	c.JSON(http.StatusOK, res)

	return

}
func CreateUser (c *gin.Context){

	createReq :=vm.CreateUserReq{}
	err1 :=c.Bind(&createReq)
	fmt.Println("creatReq:",createReq)
	if err1 != nil {
		fmt.Println("err1:",err1)
		c.JSON(http.StatusOK,err_msg.CodeMsg{Code:1,Info:"请求参数格式错误"})
		return
	}

	//if len(createReq.Id) == 0 ||len(createReq.Phone) == 0||len(createReq.Password) == 0{
	//	c.JSON(http.StatusOK,err_msg.CodeMsg{Code:1,Info:"公司名称，或账号，密码不能为空"})
	//	return
	//}
	if len(createReq.NickName) == 0 {
		c.JSON(http.StatusOK, err_msg.CodeMsg{1,"用户名字或者密码不能为空"})
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


	c.JSON(http.StatusOK,createReq)
//onlyUser :=model.User{}
//driver.MySQL().Where("Id=?",createReq.Id).Where("Phone=?",createReq.Phone).Get(&onlyUser)
//if len(onlyUser.Id)>0{
//	c.JSON(http.StatusOK,err_msg.CodeMsg{Code:1,Info:createReq.Phone+"公司账号已经存在"})
//	return
//}

}
func UserList(c *gin.Context){

}
func DeleteUser(c *gin.Context){

}




























func WebServerBase() {
	fmt.Println("This is webserver base!")

	//第一个参数为客户端发起http请求时的接口名，第二个参数是一个func，负责处理这个请求。
	http.HandleFunc("/login", loginTask)

	//服务器要监听的主机地址和端口号
	err := http.ListenAndServe("192.168.1.27:8081", nil)

	if err != nil {
		fmt.Println("ListenAndServe error: ", err.Error())
	}
}

func loginTask(w http.ResponseWriter, req *http.Request) {
	fmt.Println("loginTask is running...")

	//模拟延时
	time.Sleep(time.Second * 2)

	//获取客户端通过GET/POST方式传递的参数
	req.ParseForm()
	param_userName, found1 := req.Form["userName"]
	param_password, found2 := req.Form["password"]

	if !(found1 && found2) {
		fmt.Fprint(w, "请勿非法访问")
		return
	}

	result := NewBaseJsonBean()
	userName := param_userName[0]
	password := param_password[0]

	s := "userName:" + userName + ",password:" + password
	fmt.Println(s)

	if userName == "zhangsan" && password == "123456" {
		result.Code = 100
		result.Message = "登录成功"
	} else {
		result.Code = 101
		result.Message = "用户名或密码不正确"
	}

	//向客户端返回JSON数据
	bytes, _ := json.Marshal(result)
	fmt.Fprint(w, string(bytes))
}

type BaseJsonBean struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewBaseJsonBean() *BaseJsonBean {
	return &BaseJsonBean{}
}