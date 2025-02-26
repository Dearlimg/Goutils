package limit

import (
	"context"
	"golang.org/x/time/rate"
	"log"
	"sync"
	"time"
)

/*
有时候，test 文件中的函数并不是为了直接进行单元测试而编写的，而是作为示例代码来展示某些功能或用法。在这种情况下，可能会省略具体的测试函数。
*/

func Example() {
}

type API interface {
	ReadFile(ctx context.Context) error
	ResolveAddress(ctx context.Context) error
}

type testAPI struct {
	netWorkLimit, diskLimit, apiLimit RateLimiter //多个维度进行限制
}

func Open() API {
	apiLimit := MultiLimiter(
		rate.NewLimiter(Per(2, time.Second), 1),   //每秒的限制，防止突发请求，每 1 秒补充两个
		rate.NewLimiter(Per(10, time.Minute), 10), //每分钟的限制，设置初始池，每 10 秒补充一个
	)
	diskLimit := MultiLimiter(
		rate.NewLimiter(rate.Limit(1), 1),
	)
	netWorkLimit := MultiLimiter(
		rate.NewLimiter(Per(3, time.Second), 3),
	)
	return &testAPI{
		netWorkLimit: netWorkLimit,
		diskLimit:    diskLimit,
		apiLimit:     apiLimit,
	}
}

// ReadFile 读取文件
func (t *testAPI) ReadFile(ctx context.Context) error {
	if err := MultiLimiter(t.apiLimit, t.diskLimit).Wait(ctx); err != nil {
		return err
	}
	return nil
}

// ResolveAddress 解析地址
func (t *testAPI) ResolveAddress(ctx context.Context) error {
	if err := MultiLimiter(t.apiLimit, t.netWorkLimit).Wait(ctx); err != nil {
		return err
	}
	return nil
}

// 示例函数
func ExampleMultiLimiter() {
	defer log.Println("Done")
	apiConn := Open()
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConn.ReadFile(context.Background()); err != nil {
				log.Println("cannot read file:", err)
				return
			}
			log.Println("read file")
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConn.ResolveAddress(context.Background()); err != nil {
				log.Println("cannot resolve address:", err)
				return
			}
			log.Println("ResolveAddress")
		}()
	}
	wg.Wait()
}
