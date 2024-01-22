package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online-voice-channel/configs"
	"online-voice-channel/controller"
	"online-voice-channel/interceptor"
)

func SetupRouter() *gin.Engine {
	if configs.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	// 告诉gin框架模板文件引用的静态文件去哪里找
	//r.Static("/static", "static")
	// 告诉gin框架去哪里找模板文件
	//r.LoadHTMLGlob("templates/*")
	//r.GET("/", controller.IndexHandler)
	r.GET("/home", interceptor.ConfInterceptor(), homeHandler)

	// v1
	v1Group := r.Group("v1")
	{
		// 待办事项
		// 添加
		v1Group.POST("/todo", controller.CreateTodo)
		// 查看所有的待办事项
		v1Group.GET("/todo", controller.GetTodoList)
		// 修改某一个待办事项
		v1Group.PUT("/todo/:id", controller.UpdateATodo)
		// 删除某一个待办事项
		v1Group.DELETE("/todo/:id", controller.DeleteATodo)
	}
	loginGroup := r.Group("api")
	{
		loginGroup.POST("/login", controller.Login)
	}
	return r
}
func homeHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "success",
		"data": gin.H{"username": username},
	})
}
