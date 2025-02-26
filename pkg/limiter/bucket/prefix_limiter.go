package bucket

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"strings"
)

// PrefixLimiter 实现了 Iface 接口
type PrefixLimiter struct {
	*Limier
	*PrefixTree
}

func NewPrefixLimiter() *PrefixLimiter {
	return &PrefixLimiter{&Limier{limiterBuckets: map[string]*ratelimit.Bucket{}}, NewPrefixTree()}
}

func (p *PrefixLimiter) Key(c *gin.Context) string {
	uri := c.Request.RequestURI
	prefix := strings.Split(uri, "/")
	if result := p.Get(prefix); result != nil {
		return result.(string)
	}
	return ""
}

func (p *PrefixLimiter) testKey(uri string) string {
	prefix := strings.Split(uri, "/")
	result := p.Get(prefix)
	if result != nil {
		return result.(string)
	}
	return ""
}

func (p *PrefixLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	bucket, ok := p.limiterBuckets[key]
	return bucket, ok
}

func (p *PrefixLimiter) AddBucket(rules ...Rule) Iface {
	for _, rule := range rules {
		if _, ok := p.limiterBuckets[rule.Key]; !ok {
			//创建一个令牌桶，设置填充频率（fillInterval）、初始容量（capacity）、每秒填充的令牌数（Quantum）
			p.limiterBuckets[rule.Key] = ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Cap, rule.Quantum)
			p.Put(strings.Split(rule.Key, "/"), rule.Key)
		}
	}
	return p
}
