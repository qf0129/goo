package goo

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

type DbOption struct {
	ShowLog bool
}

func ConnectDB(opts ...DbOption) {
	var database *gorm.DB

	gormConf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// NowFunc: func() time.Time {
		// 	return time.Now().Local()
		// },
	}

	for _, opt := range opts {
		if opt.ShowLog {
			gormConf.Logger = logger.Default.LogMode(logger.Info)
		}
	}

	var dbConn gorm.Dialector
	if Config.DbEngine == "mysql" {
		dbUri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", Config.DbUser, Config.DbPsd, Config.DbHost, Config.DbPort, Config.DbDatabase)
		dbConn = mysql.Open(dbUri)
	} else if Config.DbEngine == "sqlite" {
		dbUri := fmt.Sprintf("%s.db", Config.DbDatabase)
		dbConn = sqlite.Open(dbUri)
	} else {
		logrus.Panic("InvalidDbType")
	}

	database, err := gorm.Open(dbConn, gormConf)
	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database
}

func MigrateModels(dst ...any) {
	if err := DB.AutoMigrate(dst...); err != nil {
		logrus.Panic("AutoMigrateErr:", err)
		return
	}
}
