package api

import (
	"Voichatter/models"
	"Voichatter/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
			"code": 400,
			"msg":  "账号或密码错误",
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
	username, _ := c.Get("username")
	var user models.User
	user.Username = username.(string)
	var server []models.Server
	if err := service.FindUserServersList(&user, &server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": server,
	})
}

func UserLogout(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	token := parts[1]
	if err := service.UserLogout(token); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "发生错误",
			"data": "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "退出登录成功",
		"data": "",
	})
}
