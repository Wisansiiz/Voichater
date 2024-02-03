package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online-voice-channel/models"
	"online-voice-channel/service"
)

func UserRegister(c *gin.Context) {
	var user models.User
	_ = c.ShouldBind(&user)
	if err := service.UserRegister(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

func UserLogin(c *gin.Context) {
	var user models.UserLoginResponse
	_ = c.ShouldBind(&user)
	token, err := service.UserLogin(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 200,
			"msg":  "err",
			"data": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"token": token,
		},
	})
}

func FindUserServersList(c *gin.Context) {
	var user models.User
	var server []models.Server
	_ = c.ShouldBind(&user)
	if err := service.FindUserServersList(&user, &server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": server,
	})
}
