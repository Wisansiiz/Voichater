package api

//func Login(c *gin.Context) {
//	// 1. 从请求中把数据拿出来
//	var login models.Login
//	err := c.BindJSON(&login)
//	if err != nil {
//		return
//	}
//	err = models.LoginOn(&login)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//	} else {
//		c.JSON(http.StatusOK, login)
//	}
//}
