package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"strconv"
	"time"
	"github.com/LittleCurry/marry/vm"
	"github.com/LittleCurry/misc/model"
	"github.com/LittleCurry/misc/driver"
	"github.com/LittleCurry/misc/globals"
	"regexp"
	"github.com/LittleCurry/misc/err_msg"
	"os"
	"io"
	"strings"
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

	//createReq := vm.CreateUserReq{}
	//err1 := c.Bind(&createReq)
	//fmt.Println("Phone:", createReq.Phone)
	//fmt.Println("Address:", createReq.Address)
	//fmt.Println("UserName:", createReq.UserName)
	//fmt.Println("Introduction:", createReq.Introduction)
	//fmt.Println("Birthday:", createReq.Birthday)

	user_name := c.PostForm("user_name")
	phone := c.PostForm("phone")
	passwd := c.PostForm("passwd")
	gender := c.PostForm("gender")
	birthday := c.PostForm("birthday")
	address := c.PostForm("address")
	introduction := c.PostForm("introduction")


	fileHeader, err1 := c.FormFile("head")
	strArr := strings.Split(fileHeader.Filename, ".")
	timeStr := strconv.Itoa(int(time.Now().Unix())) + "." + strArr[len(strArr)-1]
	location := "./img/" + timeStr
	if err1 != nil {
		fmt.Println("err1:", err1)
	}
	file, err2 := fileHeader.Open()
	if err2 != nil {
		fmt.Println("err2:", err2)
	}
	saveFile, err4 := os.Create(location)
	if err4 != nil {
		fmt.Println("err4:", err4)
	}
	defer saveFile.Close()

	_, err5 := io.Copy(saveFile, file)
	if err5 != nil {
		fmt.Println("err5:", err5)
	}
	chmodErr := os.Chmod(location, os.ModePerm)
	if chmodErr != nil {
		fmt.Println("chmodErr:", chmodErr)
	}



	fmt.Println("user_name:", user_name)
	fmt.Println("phone:", phone)
	fmt.Println("passwd:", passwd)
	fmt.Println("gender:", gender)
	fmt.Println("birthday:", birthday)
	fmt.Println("address:", address)
	fmt.Println("introduction:", introduction)
	fmt.Println("head:", fileHeader.Filename)


	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	if !rgx.MatchString(phone) {
		c.JSON(http.StatusOK, err_msg.CodeMsg{1,"手机号格式错误"})
		return
	}

	user := model.User{}
	user.UserId = strconv.Itoa(int(time.Now().Unix()))
	user.UserName = user_name
	user.Phone = phone
	user.Passwd = globals.MakeMd5FromString(passwd)
	user.Birthday = birthday
	user.Address = address
	user.Introduction = introduction
	user.Head = "https://media.siiva.com/img/"+timeStr
	user.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	user.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	_, err3 := driver.MySQL().Insert(&user)
	if err3 != nil {
		fmt.Println("err3:", err3)
	}

	c.JSON(http.StatusOK, user)
	return

}
func UserList(c *gin.Context) {

	users := []model.User{}
	err := driver.MySQL().Where("`del` = 0").Find(&users)
	if err != nil {
		fmt.Println("err:", err)
	}
	usersRes := []vm.UserRes{}
	copier.Copy(&usersRes, users)
	c.JSON(http.StatusOK, usersRes)

}
func DeleteUser(c *gin.Context) {

	userId := c.Query("user_id")
	fmt.Println("userId:", userId)
	if len(userId) == 0 {
		c.JSON(http.StatusOK, err_msg.CodeMsg{1, "user_id不能为空"})
		return
	}
	_, err1 := driver.MySQL().Where("user_id=?", userId).Update(&model.User{Del: 1})
	if err1 != nil {
		fmt.Println("err1:", err1)
		c.JSON(http.StatusOK, err_msg.CodeMsg{1, "删除失败"})
		return
	}
	c.JSON(http.StatusOK, err_msg.CodeMsg{0, "删除成功"})
	return
}
