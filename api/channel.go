package api

import (
	"Voichatter/models"
	"Voichatter/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateChannel(c *gin.Context) {
	serverId, _ := strconv.ParseUint(c.Param("serverId"), 10, 64)
	var channel models.Channel
	_ = c.ShouldBind(&channel)
	channel.ServerID = uint(serverId)
	userId := c.MustGet("user_id").(uint)
	if err := service.CreateChannel(&channel, userId); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":     400,
			"messages": err.Error(),
			"data":     nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"messages": "创建成功",
		"data":     nil,
	})
}
