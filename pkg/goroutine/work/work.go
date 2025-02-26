package work

import "sync"

/*
	工作池，它使用一组固定数量的工作线程来执行任务队列中的工作单元。
	可以处理：“一组耗时的任务需要执行，我们希望并发执行它们，但同时限制开发度。“的问题
*/

type Worker struct {
	taskChan chan func()  // 任务通道，用于接收待执行的任务函数
	workChan chan func()  // 工作通道
	rwMutex  sync.RWMutex // 读写锁
}

type Config struct {
	TaskChanCapacity   int // 任务队列容量
	WorkerChanCapacity int // 工作队列容量
	WorkerNum          int // 工作池数
}

func Init(config Config) *Worker {
	w := &Worker{
		taskChan: make(chan func(), config.TaskChanCapacity),
		workChan: make(chan func(), config.WorkerChanCapacity),
		rwMutex:  sync.RWMutex{},
	}
	for i := 0; i < config.WorkerNum; i++ {
		go w.work()
	}
	return w
}

func (w *Worker) work() {
	for task := range w.taskChan {
		task()
	}
}

func (w *Worker) SendTask(task func()) {
	w.rwMutex.Lock()
	defer w.rwMutex.Unlock()
	w.taskChan <- task
}
