package bucket

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

type Iface interface {
	Key(c *gin.Context) string                      //获取对应的限流器的键值对名称
	GetBucket(key string) (*ratelimit.Bucket, bool) //获取对应的限流器的键值对名称
	AddBucket(rules ...Rule) Iface                  //新增多个令牌桶规则
}

// Limier 存储令牌桶与键值对名称的映射关系
type Limier struct {
	limiterBuckets map[string]*ratelimit.Bucket
}

// Rule 令牌桶的规则
type Rule struct {
	Key          string        //自定义键值对名称
	FillInterval time.Duration //增加新桶的间隔时间
	Cap          int64         // 桶的最大容量
	Quantum      int64         // 每次到大间隔时间之后存放的桶数量
}
