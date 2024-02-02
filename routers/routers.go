package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"online-voice-channel/api"
	"online-voice-channel/configs"
	"online-voice-channel/interceptor"
)

func SetupRouter() *gin.Engine {
	if configs.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(cors.Default())
	// 静态文件去
	r.StaticFS("/static", http.Dir("./static"))

	// v1
	v1Group := r.Group("v1")
	{
		authed := v1Group.Group("/") // 需要登陆保护
		authed.Use(interceptor.ConfInterceptor())
		{
			// 添加待办事项
			authed.POST("/todo", api.CreateTodo)
			// 查看所有的待办事项
			authed.GET("/todo", api.GetTodoList)
			// 修改某一个待办事项
			authed.PUT("/todo/:id", api.UpdateATodo)
			// 删除某一个待办事项
			authed.DELETE("/todo/:id", api.DeleteATodo)
		}
	}
	v := r.Group("api")
	{
		v.POST("/register", api.UserRegister)
		v.POST("/login", api.UserLogin)
		v.GET("/servers-list", api.FindUserServersList)
		v.GET("/ws", api.Ws)
		v.GET("/history", api.FindMessage)
	}
	return r
}
