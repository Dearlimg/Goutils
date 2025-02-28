package task

import (
	"context"
	"github.com/Dearlimg/Goutils/pkg/goroutine/heal"

	"log"
	"time"
)

type Task struct {
	Name            string          // 任务名
	Ctx             context.Context // 上游 ctx
	TaskDuration    time.Duration   // 任务执行周期
	TimeoutDuration time.Duration   // 超时时长
	F               DoFunc          // 执行程序
}

type DoFunc func(parentCtx context.Context)

// NewTickerTask 创建一个被管理的定时任务，并返回可以监听的管理者的心跳
// 定时任务通过发送心跳信号来告知管理者自己仍然处于活跃状态。
// 管理者收到心跳信号后，可以根据心跳的频率和稳定性来监控任务的健康状况。
// 如果管理者在一段时间内未收到任务的心跳信号，可能会认为任务执行者已经出现了故障或不可用，并采取相应的措施，比如重新分配任务或触发报警。
func NewTickerTask(task Task) <-chan struct{} {
	startFun := func(ctx context.Context, pulseInterval time.Duration) <-chan struct{} {
		ticker := time.NewTicker(task.TaskDuration) // 定时任务
		pulse := time.NewTicker(pulseInterval)      // 定期心跳
		heartBeat := make(chan struct{})
		go func() {
			defer ticker.Stop() // 关闭后停止定时器
			defer pulse.Stop()  // 关闭后停止回复心跳
			now := time.Now()
			task.F(ctx) // 最少会执行一次
			log.Println("first exec task:", task.Name, "cost time:", time.Since(now))
			for {
				select {
				case <-ticker.C: // 通知计时器必须等多长时间
					now = time.Now()
					task.F(ctx)
					log.Println("task: try to exec task:", task.Name, "cost time:", time.Since(now))
				case <-pulse.C: // 等多长时间后必须回复，以提醒管理者这个定时任务是正常运行的
					heartBeat <- struct{}{}
				case <-ctx.Done():
					log.Println("task: over by stewart:", task.Name)
					return
				}
			}
		}()
		return heartBeat
	}
	// 调用之前定义的 NewSteward 函数，并传递了相应的参数 task.Name、task.TimeoutDuration 和 startFun，返回了一个函数。
	// 然后，立即调用返回的函数，并传递了参数 task.Ctx 和 task.TimeoutDuration。
	return heal.NewSteward(task.Name, task.TimeoutDuration, startFun)(task.Ctx, task.TimeoutDuration)
}
