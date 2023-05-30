package global

import (
	"github.com/sony/sonyflake"
	"strconv"
	"sync"
)

var idInstance *sonyflake.Sonyflake

var once = new(sync.Once)

func init() {
	once.Do(func() {
		settings := sonyflake.Settings{}
		idInstance = sonyflake.NewSonyflake(settings) // 用配置生成sonyflake节点
	})
}

func ID() string {
	v, _ := idInstance.NextID()
	return strconv.FormatUint(v, 10)
}
