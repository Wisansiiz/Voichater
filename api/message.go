package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online-voice-channel/dao"
	"online-voice-channel/models"
)

var msg []models.Message

func Ws(c *gin.Context) {
	// WebSocket处理程序
	models.HandleWebSocket(c.Writer, c.Request, dao.DB)
}

func FindMessage(c *gin.Context) {
	_ = c.ShouldBind(&msg)
	if err := models.FindHistory(&msg, c.Request, dao.DB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "发生错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": msg,
	})
}
