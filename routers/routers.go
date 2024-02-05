package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online-voice-channel/api"
	"online-voice-channel/configs"
	"online-voice-channel/interceptor"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	if configs.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(Cors())
	// 静态文件去
	r.StaticFS("/static", http.Dir("./static"))

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
