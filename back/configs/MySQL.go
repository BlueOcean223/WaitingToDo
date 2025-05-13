package configs

import (
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

	return nil
}
