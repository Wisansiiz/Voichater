package api

import (
	"Voichatter/models"
	"Voichatter/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var messages []models.Message

func Ws(c *gin.Context) {
	// WebSocket处理程序
	service.HandleWebSocket(c)
}

func FindMessage(c *gin.Context) {
	_ = c.ShouldBind(&messages)
	channelID := c.Query("channelID")
	if err := service.FindHistory(&messages, channelID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":     400,
			"messages": "发生错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"messages": "success",
		"data":     messages,
	})
}
