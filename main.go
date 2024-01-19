package main

import (
	"online-voice-channel/dao"
	"online-voice-channel/models"
	"online-voice-channel/routers"
	"online-voice-channel/setting"

	"fmt"
	"os"
)

const defaultConfFile = "./conf/config.yaml"

func main() {
	confFile, err := os.ReadFile(defaultConfFile)
	if err != nil {
		panic(err)
	}
	// 加载配置文件
	if err := setting.Init(confFile); err != nil {
		fmt.Printf("load config from file failed, err:%v\n", err)
		return
	}
	// 创建数据库
	// 连接数据库
	err = dao.InitMySQL(setting.Conf.DatabaseConfig)
	if err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer dao.Close(dao.DB) // 程序退出关闭数据库连接
	// 模型绑定
	err = dao.DB.AutoMigrate(&models.Todo{})
	if err != nil {
		return
	}
	// 注册路由
	r := routers.SetupRouter()
	if err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port)); err != nil {
		fmt.Printf("server startup failed, err:%v\n", err)
	}
}
