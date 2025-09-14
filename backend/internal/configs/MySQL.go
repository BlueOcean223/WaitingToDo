package configs

import (
	"backend/internal/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConfig struct {
	Dsn string `yaml:"dsn"`
	// 连接池配置
	MaxIdleConns    int `yaml:"max_idle_conns"`     // 最大空闲连接数
	MaxOpenConns    int `yaml:"max_open_conns"`     // 最大打开连接数
	ConnMaxLifetime int `yaml:"conn_max_lifetime"`  // 连接最大生命周期
	ConnMaxIdleTime int `yaml:"conn_max_idle_time"` // 连接最大空闲时间
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

	// 连接池配置
	sqlDB, err := MysqlDb.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(AppConfigs.MySQLConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfigs.MySQLConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(AppConfigs.MySQLConfig.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(AppConfigs.MySQLConfig.ConnMaxIdleTime) * time.Second)

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
		&models.InviteCode{},
	)
}
