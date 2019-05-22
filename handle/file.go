package handle

import (
	"github.com/gin-gonic/gin"
	"fmt"
)

func GetImage(c *gin.Context) {
	imgName := c.Param("name")
	fmt.Println("img_name:", imgName)
	c.File("./img/" + imgName)
	return
}
