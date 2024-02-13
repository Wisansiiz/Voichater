package api

import (
	"Voichatter/models"
	"Voichatter/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var msg []models.Message

func Ws(c *gin.Context) {
	// WebSocket处理程序
	service.HandleWebSocket(c)
}

func FindMessage(c *gin.Context) {
	_ = c.ShouldBind(&msg)
	channelID := c.Query("channelID")
	if err := service.FindHistory(&msg, channelID); err != nil {
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
