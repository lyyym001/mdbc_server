package models_sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"mdbc_server/internal/config"
	"mdbc_server/lframework/zlog"
	"time"
)

var DB *gorm.DB
var DirVersion int64

func NewDB() {
	//dsn := string(config.YamlConfig.Mysql.Dns)
	//"root:root@tcp(127.0.0.1:3306)/mspace?charset=utf8mb4&parseTime=True&loc=Local
	//"root:root@tcp(127.0.0.1:3306)/mspace?charset=utf8mb4&parseTime=True&loc=Local"
	//fmt.Println("config.YamlConfig.Sqlite.Dns = ", config.YamlConfig.Sqlite.Dns)
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags),
	//	logger.Config{
	//		SlowThreshold: time.Second, //慢SQL阈值
	//		LogLevel:      logger.Info, //级别
	//		Colorful:      true,        //彩色
	//	},
	//)
	db, err := gorm.Open(sqlite.Open(config.YamlConfig.Sqlite.Dns), &gorm.Config{
		//Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	sqldb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(100)
	sqldb.SetConnMaxLifetime(time.Hour)

	zlog.Debugf("connect db success - dbName = %s", config.YamlConfig.Sqlite.Dns)
	db.AutoMigrate(&DeviceBasic{}, &DirBasic{}, &CourseBasic{}, &ShoucangBasic{}, &RecordBasic{}, &WorkBasic{}, &WorkRecordBasic{}, &SysBasic{}) //&RoomBasic{}, &RoomUser{},
	DB = db

	// 2.默认创建表
	var count int64
	err = DB.Model(SysBasic{}).Where("sid = ?", 1).Count(&count).Error
	if err != nil || count == 0 {
		u := &SysBasic{DirVersion: 0, Sid: 1}
		DB.Create(&u)
	}

	// 2.读取系统状态
	sData := &SysBasic{}
	err = DB.Where("sid = ?", 1).First(sData).Error
	if err == nil {
		//设置状态
		DirVersion = sData.DirVersion
	}
}
