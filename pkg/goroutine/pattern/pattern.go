package pattern

import (
	"context"
	"sync"
)

// Or 监听多个 ctx，只要有一个返回消息就返回
// 函数接受一个或多个对象作为参数。
// 它会返回一个新的上下文对象，该上下文对象会在传入的多个上下文对象中的任何一个完成时完成。
// 这意味着只要其中一个上下文被取消或超时，新的上下文就会被取消或超时。
func Or(ctx ...context.Context) context.Context {
	switch len(ctx) {
	case 0:
		return nil
	case 1:
		return ctx[0]
	}
	orCtx, cancel := context.WithCancel(context.Background())
	// 启动一个新的 goroutine 来处理上下文的取消情况
	go func() {
		// 在函数退出的时候，调用取消函数，确保即使释放资源
		defer cancel()
		switch len(ctx) {
		case 2:
			select {
			case <-ctx[0].Done():
			case <-ctx[1].Done():
			}
		default:
			select {
			case <-ctx[0].Done():
			case <-ctx[1].Done():
			case <-ctx[2].Done():
			case <-Or(append(ctx[3:], orCtx)...).Done(): // 如果有更多的上下文，通过递归调用 Or 函数处理剩余的上下文，并监听返回的上下文的取消情况
			}
		}
	}()
	return orCtx // 返回新创建的上下文对象
}

// Bridge 将一个通道的通道转换为一个单独的通道，并且保持数据的持续传输（这个是按顺序读完一个 channel 才会选择下一个 channel）
// 参数：ctx：上下文，用于跟踪函数的执行状态和控制函数的生命周期
// 参数：chanStream：接收一个元素类型为 <-chan interface{} 的只接收通道的通道
// 返回值：返回一个元素类型为 interface{} 的只发送通道
func Bridge(ctx context.Context, chanStream <-chan <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			var stream <-chan interface{}
			select {
			case chanS, ok := <-chanStream: // 读取 chanStream 中的 channel
				if !ok {
					return
				}
				stream = chanS
			case <-ctx.Done():
				return
			}
			for val := range OrDone(ctx, stream) { // 读取 channel 中的内容发送回去
				select {
				case <-ctx.Done():
				case valStream <- val:
				}
			}
		}
	}()
	return valStream
}

// OrDone 安全地读取通道 c 中的数据
func OrDone(ctx context.Context, c <-chan interface{}) <-chan interface{} {
	varStream := make(chan interface{})
	go func() {
		defer close(varStream)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case varStream <- v: // 将从通道 c 中读取的数据发送到 valStream 中
				case <-ctx.Done():
				}
			}
		}
	}()
	return varStream
}

// Tee 读取 in 数据并同时发送两个接收的 channel
func Tee(ctx context.Context, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for v := range OrDone(ctx, in) {
			var out1, out2 = out1, out2 // 本地版本，隐藏外界变量
			for i := 0; i < 2; i++ {
				select {
				case <-ctx.Done():
					return
				case out1 <- v:
					out1 = nil // 同时写入后关闭副本 channel 来阻塞防止二次写入
				case out2 <- v:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

// FanIn 从多个 channel 中合并数据到一个 channel
func FanIn(ctx context.Context, channels []<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})
	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-ctx.Done():
			case multiplexedStream <- i:
			}
		}
	}
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}
	// 不在主函数中使用 wg.Wait() 阻塞，而是新开一个 goroutine，
	// 如此主函数就不必等所有的 goroutine 都执行完毕再返回结果通道，而是立刻返回通道，以便可以立刻读取数据
	// 但是会在后台进行通道数据的输入与通道的关闭
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()
	return multiplexedStream
}

// Take 去除 num 个数后结束
func Take(ctx context.Context, valueStream <-chan interface{}, num int) <-chan interface{} {
	results := make(chan interface{})
	go func() {
		defer close(results)
		for i := 0; i < num; i++ {
			select {
			case <-ctx.Done():
				return
			case results <- valueStream:
			}
		}
	}()
	return results
}

// RepeatFn 重复调用函数(返回一个数据值为 fn() 函数的 channel)
func RepeatFn(ctx context.Context, fn func() interface{}) <-chan interface{} {
	results := make(chan interface{})
	go func() {
		defer close(results)
		for {
			select {
			case <-ctx.Done():
				return
			case results <- fn():
			}
		}
	}()
	return results
}

// Repeat 重复生成值
func Repeat(ctx context.Context, values ...interface{}) chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-ctx.Done():
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}
