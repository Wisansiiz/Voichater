package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"online-voice-channel/configs"
	"time"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB https://gorm.io/zh_CN/docs/
var DB *gorm.DB

func InitMySQL(cfg *configs.MySqlConfig) {
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据版本自动配置
		}), &gorm.Config{
			Logger: ormLogger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	if err != nil {
		fmt.Println("连接失败")
		panic(err)
	}
	sqlDB, _ := db.DB()

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // SetMaxIdleConns 设置空闲连接池中的最大连接数
	sqlDB.SetMaxOpenConns(100)          // SetMaxOpenConns 设置数据库的最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // SetConnMaxLifetime 设置连接可重用的最大时间量
	DB = db                             // 设置全局DB
}
func Close(DB *gorm.DB) {
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("sqlDB错误")
		panic(err)
	}
	err = sqlDB.Close()
	if err != nil {
		fmt.Println("关闭失败")
		panic(err)
	}
} // 关闭链接
