package limit

import (
	"context"
	"golang.org/x/time/rate"
	"sort"
	"time"
)

// RateLimiter 限流接口
type RateLimiter interface {
	Wait(ctx context.Context) error //阻塞等待
	Limit() rate.Limit
}

type multiLimiter struct {
	limiters []RateLimiter
}

// MultiLimiter 聚合多个 RateLimiter，并将速率由小到大排序
func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

// Wait 会循环遍历多层限速器 multiLimiter 中所有的限速器并索要令牌，只有当所有的限速器规则都满足后，才会正常执行后续的操作
func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Limit 返回当前限制速率
func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

// Per 返回速率为每 duration，eventCount 个请求
func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}
