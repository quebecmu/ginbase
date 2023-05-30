package setting

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xxandjg/ginbase/global"
	"github.com/xxandjg/ginbase/http/middleware"
	"go.uber.org/zap"
)

func ServerInit() {
	r := gin.New()
	r.Use(middleware.Cors())
	r.Use(global.GinLogger(), global.GinRecovery(true))
	//router.Init(r)
	//routes := make([]string, 0)
	//r.Routes()
	//for _, route := range r.Routes() {
	//	routes = append(routes, route.Method+" "+route.Path)
	//}
	if err := r.Run(fmt.Sprintf(":%d", App.Application.Port)); err != nil {
		zap.L().Sugar().Errorf("server run failed, err:%v\n", err)
		return
	}
}
