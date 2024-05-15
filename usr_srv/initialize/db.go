package initialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mxshop_srvs/usr_srv/global"
	"mxshop_srvs/usr_srv/model"
	"os"
	"time"
)

func InitDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:        time.Second, // Slow SQL threshold
			LogLevel:             logger.Info, // Log level
			ParameterizedQueries: false,
			Colorful:             true, // color
		},
	)

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	//dsn := "root:123456@tcp(192.168.57.128:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlInfo := global.ServerConfig.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlInfo.User, mysqlInfo.Password, mysqlInfo.Host, mysqlInfo.Port, mysqlInfo.Name)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})

	if err != nil {
		panic("连接数据库失败\n" + err.Error())
	}

	global.DB.AutoMigrate(&model.User{})
}
