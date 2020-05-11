package datasourse

import (
	"E-commerce/common/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var dbInstance *gorm.DB

// 单例模式
func GetMysqlInstance() *gorm.DB {
	if dbInstance == nil {
		dbInstance = initDB()
	}

	return dbInstance
}

func initDB() *gorm.DB {
	db, err := gorm.Open(conf.DataBaseSetting.Type, conf.DataBaseSetting.SqlDSN)
	// 设置为true,数据操作日志可以数据在控制台
	db.LogMode(true)
	// 默认情况下创建表明结构为复数
	db.SingularTable(true)

	// Error
	if err != nil {
		conf.AppSetting.Logger.Error("连接数据库不成功", err, conf.DataBaseSetting)
	} else {
		conf.AppSetting.Logger.Info("数据库 链接成功\n")
	}

	//设置连接池
	//空闲
	db.DB().SetMaxIdleConns(conf.DataBaseSetting.MaxIdle)
	//打开
	db.DB().SetMaxOpenConns(conf.DataBaseSetting.MaxOpen)
	//超时
	db.DB().SetConnMaxLifetime(time.Second * conf.DataBaseSetting.IdleTimeout)

	return db
}
