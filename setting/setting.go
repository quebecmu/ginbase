package setting

import (
	"github.com/gin-gonic/gin"
	"github.com/xxandjg/ginbase/global"
	"go.uber.org/zap"
	"sync"
)

func Init() {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 加载配置文件
		if err := ConfigInit(gin.Mode()); err != nil {
			zap.L().Sugar().Errorf("load config failed, err:%v\n", err)
		}
	}()
	wg.Wait()
	//fmt.Println(App.MysqlInfo)
	//fmt.Println(App.LogInfo)
	err := global.InitLogger(App.LogInfo.MaxBackups, App.LogInfo.MaxAge, App.LogInfo.MaxSize, App.LogInfo.Path, App.LogInfo.Level)
	if err != nil {
		return
	}
	ServerInit()
}
