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
	const defaultConfFile = "configs/locales/config.yaml"
	configs.InitConfig(defaultConfFile)
	// 初始化mysql
	dao.InitMySQL(configs.Conf.MySql)
	// 初始化redis
	dao.InitRedis(configs.Conf.Redis)
	defer dao.Close(dao.DB) // 程序退出关闭数据库连接
	// 模型绑定
	err := dao.DB.AutoMigrate(&models.User{}, &models.Message{},
		&models.Channel{}, &models.Server{}, &models.Member{},
	)
	if err != nil {
		panic(err)
	}
	// 注册路由
	r := routers.SetupRouter()
	if err = r.Run(fmt.Sprintf(":%d", configs.Conf.Port)); err != nil {
		fmt.Printf("server startup failed, err:%v\n", err)
		panic(err)
	}
}
