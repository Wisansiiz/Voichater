package routers

import (
	"Voichatter/api"
	"Voichatter/configs"
	"Voichatter/interceptor"
	"Voichatter/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	if configs.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(middleware.Cors())
	v := r.Group("api")
	{
		v.POST("/register", api.UserRegister)
		v.POST("/login", api.UserLogin)
		v.GET("/ws", api.Ws)
		authed := v.Group("/") // 需要登陆保护
		authed.Use(interceptor.ConfInterceptor())
		{
			authed.POST("/auth", interceptor.Auth)
			authed.POST("/logout", api.UserLogout)
			authed.GET("/servers-list", api.FindUserServersList)
			authed.GET("/history", api.FindMessage)
		}
	}
	return r
}
