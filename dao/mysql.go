package dao

import (
	"online-voice-channel/setting"
	"time"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB https://gorm.io/zh_CN/docs/
var DB *gorm.DB

func InitMySQL(cfg *setting.DatabaseConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接失败")
	}
	sqlDB, err := DB.DB()

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // SetMaxIdleConns 设置空闲连接池中的最大连接数
	sqlDB.SetMaxOpenConns(100)          // SetMaxOpenConns 设置数据库的最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // SetConnMaxLifetime 设置连接可重用的最大时间量

	return err
}
func Close(DB *gorm.DB) {
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("sqlDB错误")
	}
	err = sqlDB.Close()
	if err != nil {
		fmt.Println("关闭失败")
	}
} // 关闭链接
