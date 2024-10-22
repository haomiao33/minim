package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var Db *gorm.DB

func Init(user string, password string, host string, port int, database string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s"+
		"?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		database)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,         // Don't include params in the SQL log
			Colorful:                  true,         // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	// 设置数据库连接池空闲连接数
	dbInstance, err := db.DB()
	if err != nil {
		log.Fatal("failed to open database")
	}
	// 打开连接
	dbInstance.SetMaxIdleConns(2)
	Db = db
}
