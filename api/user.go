package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online-voice-channel/models"
	"online-voice-channel/pkg/utils/translator"
	"online-voice-channel/service"
)

func UserRegister(c *gin.Context) {
	var user models.User
	if err := translator.ReErr(user); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
	}
	_ = c.ShouldBind(&user)
	if err := service.UserRegister(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

func FindUserServersList(c *gin.Context) {
	var user models.User
	var server []models.Server
	if err := translator.ReErr(user); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
	}
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
