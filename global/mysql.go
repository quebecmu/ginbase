package global

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	MysqlDB = new(gorm.DB)
)

func InitMysql(url string) Error {
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		zap.L().Sugar().Error(err)
	}
	MysqlDB = db
	zap.L().Sugar().Info("init MysqlDB success")
	return SUCCESS
}
