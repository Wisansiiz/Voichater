package main

import (
	"fmt"
	"online-voice-channel/configs"
	"online-voice-channel/dao"
	"online-voice-channel/models"
	"online-voice-channel/routers"
)

func main() {
	// 加载配置文件
	err := configs.InitConfig()
	if err != nil {
		panic(err)
	}
	// 连接数据库
	err = dao.InitMySQL(configs.Conf.MySql)
	if err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		panic(err)
	}
	defer dao.Close(dao.DB) // 程序退出关闭数据库连接
	// 模型绑定
	err = dao.DB.AutoMigrate(&models.Todo{})
	if err != nil {
		return
	}
	// 注册路由
	r := routers.SetupRouter()
	if err := r.Run(fmt.Sprintf(":%d", configs.Conf.Port)); err != nil {
		fmt.Printf("server startup failed, err:%v\n", err)
		return
	}
}
