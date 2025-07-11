package configs

import (
	"back/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConfig struct {
	Dsn string `yaml:"dsn"`
}

var MysqlDb *gorm.DB

func InitMysqlConnection() error {
	dsn := AppConfigs.MySQLConfig.Dsn
	var err error

	// 连接数据库
	MysqlDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 当数据库不存在相应的表结构时，自动创建
	err = InitMysql(MysqlDb)
	if err != nil {
		return err
	}

	return nil
}

// InitMysql 初始化数据库表结构
func InitMysql(mysqlDb *gorm.DB) error {
	return mysqlDb.AutoMigrate(
		&models.User{},
		&models.Message{},
		&models.Friend{},
		&models.Task{},
		&models.TeamTask{},
		&models.File{},
		&models.Image{},
		&models.TaskNoticeHistory{},
	)
}
