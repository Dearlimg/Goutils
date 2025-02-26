package singleflight

import "sync"

/*
防止缓存击穿
*/

/*
缓存雪崩：缓存在同一时刻全部失效，造成瞬时DB请求量大、压力骤增，引起雪崩。缓存雪崩通常因为缓存服务器宕机、缓存的 key 设置了相同的过期时间等引起。
缓存击穿：一个存在的key，在缓存过期的一刻，同时有大量的请求，这些请求都会击穿到 DB ，造成瞬时DB请求量大、压力骤增。
缓存穿透：查询一个不存在的数据，因为不存在则不会写到缓存中，所以每次都会去请求 DB，如果瞬间流量过大，穿透到 DB，导致宕机。
*/

// call 代表正在进行中 或已经结束的请求。使用 sync.WaitGroup 锁避免重入
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group 是 singleflight 的注数据结构，管理不同 key 的请求（call）
type Group struct {
	mu sync.Mutex
	m  map[interface{}]*call
}

func NewGroup() *Group {
	return &Group{m: make(map[interface{}]*call), mu: sync.Mutex{}}
}

// Do 保证 key 所对应的 fn 函数同一时刻只会执行一次
// 如果并发协程之间不需要消息传递，非常适合 sync.WaitGroup
func (g *Group) Do(key interface{}, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock() //加锁，保证并发安全
	if g.m == nil {
		g.m = make(map[interface{}]*call)
	}
	//如果发现有函数正在运行，则等待其运行并返回其返回值
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	// 唯一一个运行的函数，添加到 map 中，并且设置锁
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	c.val, c.err = fn()
	//标记着函数执行完成
	c.wg.Done()
	//完成任务后从 map 中删除
	g.mu.Lock() //加锁，保证并发安全
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err
}
