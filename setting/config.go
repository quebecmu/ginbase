package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/xxandjg/ginbase/entity"
	"github.com/xxandjg/ginbase/global"
	"go.uber.org/zap"
)

var App = new(entity.System)

// ConfigInit 初始化并监听配置文件
func ConfigInit(m string) error {
	viper.SetConfigFile("./config/app-" + m + ".yaml")

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err2 := viper.Unmarshal(&App); err2 != nil {
			zap.L().Sugar().Errorf("unmarshal to Conf failed, err:%v", err2)
		}
		if err := global.InitMysql(App.MysqlInfo.Url); err.GetCode() != 10000 {
			zap.L().Sugar().Error(err)
		}
		if err := global.InitRedis(App.RedisInfo); err.GetCode() != 10000 {
			zap.L().Sugar().Error(err)
		}
		//if err := global.RulesInit(); err.GetCode() != 20000 {
		//	zap.L().Sugar().Error(err)
		//}
	})
	err := viper.ReadInConfig()
	if err != nil {
		zap.L().Sugar().Errorf("unmarshal to Conf failed, err:%v", err)
	}
	if err := viper.Unmarshal(&App); err != nil {
		zap.L().Sugar().Errorf("unmarshal to Conf failed, err:%v", err)
	}

	if err := global.InitMysql(App.MysqlInfo.Url); err.GetCode() != 10000 {
		zap.L().Sugar().Error(err)
	}
	if err := global.InitRedis(App.RedisInfo); err.GetCode() != 10000 {
		zap.L().Sugar().Error(err)
	}
	return err
}
