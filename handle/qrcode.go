package handle

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/err_msg"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"os"
	"gitlab.com/SiivaVideoStudio/cloud_server/misc/globals"
)

func GetQrcode(c *gin.Context) {
	taskId, _ := c.GetQuery("taskId")
	fmt.Println("taskId:", taskId)
	if len(taskId) <= 0 || taskId == "undefined" || len(taskId) > 32 {
		c.JSON(http.StatusBadRequest, err_msg.CodeMsg{1, "taskId格式错误"})
		return
	}

	fileName, _ := requestQrcode("pages/main_page/main_page",  taskId)
	//c.JSON(http.StatusOK, struct {
	//	Img string `json:"img"`
	//}{"https://iva.siiva.com/me_photo/qrcode_img/" + fileName})
	c.JSON(http.StatusOK, struct {
		Img string `json:"img"`
	}{"https://siiva-video.oss-cn-hangzhou.aliyuncs.com/qrcode/" + fileName})
	return
}

func requestQrcode(page string, taskId string) (string, error) {
	urlTemplate := "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
	accessToken, tokenErr := getAccess("wxb4156cd7e9c16be9", "fa5044fa0c707f508c7f6522b342b4e7")

	if len(accessToken) <= 0 {
		return "", tokenErr
	}

	requestUrl := fmt.Sprintf(urlTemplate, accessToken)

	param := struct {
		Page  string `json:"page"`
		Scene string `json:"scene"`
	}{page,  taskId}
	paramJson, _ := json.Marshal(param)

	res, err := http.Post(requestUrl, "application/json;charset=utf-8", bytes.NewBuffer(paramJson))
	//fmt.Println("res:", res)
	if err != nil {
		fmt.Println("err0:", err)
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println("msg0:", string(body))
	}
	result, err := ioutil.ReadAll(res.Body)
	//fmt.Println("result:", string(result))
	res.Body.Close()
	if err != nil {
		fmt.Println("err1:", err)
		fmt.Println("msg1:", string(result))
	}
	fileName := writePng(result, taskId)
	return fileName, nil
}

func getAccess(appid string, secret string) (string, error) {

	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appid + "&secret=" + secret
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err2:", err)
		return "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	type AccessTokenRes struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	c := &AccessTokenRes{}
	err = json.Unmarshal(body, c)
	if err != nil {
		fmt.Println("err3:", err)
		return "", err
	}
	return c.AccessToken, nil

}

func writePng(content []byte, taskId string) string {
	//filename := string(globals.Krand(20, globals.KC_RAND_KIND_LOWER))
	filename := taskId + ".png"
	f, err := os.Create(filename)

	defer f.Close()
	if err != nil {
		fmt.Println("err3:", err)
	}

	_, err = f.Write(content)
	if err != nil {
		fmt.Println("err4:", err)
	}

	uploadErr := globals.PutFile("./"+filename, "qrcode/"+taskId+".png")
	if uploadErr == nil {
		// 删除本地文件
		//fmt.Println("请删除")
		//err2 := os.Remove(filename)
		//if err2 != nil {
		//	fmt.Println("删除文件err:", err2)
		//}
	}


	return filename
}

func QrcodeImg(c *gin.Context) {
	imgName := c.Param("name")
	fmt.Println("img_name:", imgName)
	c.File("./" + imgName)
	return
}