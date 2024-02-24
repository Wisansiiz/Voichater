package api

import (
	"Voichatter/models"
	"Voichatter/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

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
