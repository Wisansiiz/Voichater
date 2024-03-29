package interceptor

import (
	"Voichatter/dao"
	"Voichatter/pkg/utils/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ConfInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		token := parts[1]
		result, err := dao.RedisClient.Exists(dao.RedisContext, token).Result()
		if err != nil {
			return
		}
		if result > 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": 2006,
				"msg":  "Token在黑名单中",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("user_id", mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("ExpiresAt")来获取当前请求的用户信息}
	}
}
