package routers

import (
	"Voichatter/api"
	"Voichatter/configs"
	"Voichatter/interceptor"
	"Voichatter/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	if configs.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(middleware.Cors())

	v := r.Group("api")
	v.GET("/yy", api.HandleWebSocket)
	v.GET("/ws", api.Ws)
	{
		v.POST("/register", api.UserRegister)
		v.POST("/login", api.UserLogin)
		authed := v.Group("/") // 需要登陆保护
		authed.Use(interceptor.ConfInterceptor())
		{
			authed.GET("/auth", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code":     200,
					"messages": "欢迎回来",
					"data":     c.MustGet("user_id"),
				})
			})

			//authed.GET("/ws", api.Ws)
			authed.POST("/logout", api.UserLogout)
			authed.GET("/servers-list", api.FindUserServersList)
			authed.GET("/history", api.FindMessage)
			authed.POST("/create-server", api.CreateServer)
			authed.POST("/join-server", api.JoinServer)
			//authed.GET("/get-server-name", api.FindServerName)
		}
	}
	return r
}
