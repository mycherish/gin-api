package models

import (
	"fmt"
	"go-api/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func init() {
	// 使用 config 必须初始化
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	// dsn := "root:@tcp(127.0.0.1:3306)/test_go-api?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		config.C.DB.DbUsername,
		config.C.DB.DbPassword,
		config.C.DB.DbHost,
		config.C.DB.DbPort,
		config.C.DB.DbDatabase,
		config.C.DB.DbCharset,
		config.C.DB.DbParseTime,
		config.C.DB.DbLoc)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("连接数据库失败, error: ", err)
	}
}
