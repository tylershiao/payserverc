package mysql

import (
	"server/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MysqlClient sql连接对象
type MysqlClient struct {
	client *gorm.DB
}

func InitMysqlClient(dsn string, debug bool) (*MysqlClient, error) {
	dsn += "?charset=utf8mb4&parseTime=True&loc=Local"
	// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	var ormLogger logger.Interface
	if debug {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}

	user_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: ormLogger,
	})
	if err != nil {
		return nil, err
	}

	if config.GetBase().Release {
		sqlDB, err := user_db.DB()
		if err != nil {
			return nil, err
		}
		if config.GetBase().Mysql.MaxIdleConns > 0 {
			sqlDB.SetMaxIdleConns(config.GetBase().Mysql.MaxIdleConns)
		}
		if config.GetBase().Mysql.ConnMaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(config.GetBase().Mysql.ConnMaxLifetime) * time.Minute)
		}
	}

	InitBatchQueue(user_db)

	return &MysqlClient{client: user_db}, nil
}

func (m *MysqlClient) C() *gorm.DB {
	return m.client
}

// Begin 开启事务
func (u *MysqlClient) Begin() *gorm.DB {
	return u.client.Begin()
}
