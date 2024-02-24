package api

import (
	"Voichatter/models"
	"Voichatter/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func FindServerName(c *gin.Context) {
	userId, _ := c.Get("user_id")
	serverId, err := strconv.ParseUint(c.Query("server_id"), 10, 64) // 10是基数（十进制），64表示位数
	sn, err := service.FindServerName(userId.(uint), uint(serverId))
	if err != nil {
		c.JSON(200, gin.H{
			"code":     400,
			"messages": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"server_name": sn,
		},
	})
}

func GetServerMembers(c *gin.Context) {
	serverId, _ := strconv.ParseUint(c.Param("serverId"), 10, 64)
	var users []models.UserList4Server
	err := service.GetServerMembers(uint(serverId), &users)
	if err != nil {
		c.JSON(200, gin.H{
			"code":     400,
			"messages": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"members": users,
		},
	})
}
