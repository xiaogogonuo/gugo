package gugo

import (
	"context"
	"fmt"
	"time"
)

const (
	MaxIdle   = 10          // 默认最大休眠次数
	HeartBeat = time.Second // 默认心跳检测间隔时间
)

var ctx, cancel = context.WithCancel(context.Background())

type engine struct {
	maxIdle   uint64        // 最大休眠次数
	drag      chan uint64   // 阻尼器
	heartbeat time.Duration // 心跳检测间隔时间
	*spider
	*pipeline
	*scheduler
	*downloader
}

func newEngine() *engine {
	return &engine{
		maxIdle:    MaxIdle,
		heartbeat:  HeartBeat,
		drag:       make(chan uint64, 1),
		spider:     newSpider(),
		pipeline:   newPipeline(),
		scheduler:  newScheduler(),
		downloader: newDownloader(),
	}
}

// coordinate 引擎协调各组件工作
func (e *engine) coordinate() {
	e.monitorEngine()
	e.monitorPipeline()
	e.roundRobin()
	<-e.drag
	e.summary()
}

// roundRobin 轮询调度
func (e *engine) roundRobin() {
	go func() {
		for {
			select {
			case req := <-e.reqBuf:
				e.concurrentRequest <- struct{}{}
				go e.download(req, e.concurrentRequest, e.reqBuf, e.resBuf)
			case res := <-e.resBuf:
				e.concurrentResponse <- struct{}{}
				go e.response(res)
			case <-ctx.Done():
				return
			}
		}
	}()
}

// monitor 引擎监控单元，监控项：scheduler、downloader、spider
func (e *engine) monitorEngine() {
	go func() {
		var count uint64
		defer func() {
			cancel()
			e.drag <- count
		}()
		for {
			if e.idle() {
				count++
			}
			if count >= e.maxIdle {
				// 再次检查调度器的空闲状态，确保它已经可以被停止。
				if e.idle() {
					break
				} else {
					// 如果发现调度器没有休眠，则重置休眠计数
					if count > 0 {
						count = 0
					}
				}
			}
			time.Sleep(e.heartbeat)
			count++
		}
	}()
}

// monitorPipeline 管道监控单元，监控项：pipeline
func (e *engine) monitorPipeline() {
	go func() {
		var count uint64
		for {
			if e.Empty() {
				count++
			}
			if count >= e.maxIdle {
				if e.Empty() {
					close(e.pipeBuf)
					return
				} else {
					count = 0
				}
			}
			time.Sleep(e.heartbeat)
		}
	}()
}

// idle 引擎休眠逻辑
func (e *engine) idle() bool {
	if e.scheduler.HandlingNumber() == 0 &&
		e.downloader.HandlingNumber() == 0 &&
		e.spider.HandlingNumber() == 0 {
		return true
	}
	return false
}

// SetMaxIdle 设置最大休眠次数
func (e *engine) SetMaxIdle(n uint64) {
	e.maxIdle = n
}

// SetHearBeat 设置心跳检测间隔时间
func (e *engine) SetHearBeat(heartbeat time.Duration) {
	e.heartbeat = heartbeat
}

// 统计信息
func (e *engine) summary() {
	fmt.Println()
	fmt.Println("* * * * * * * * * * * * * * * * 统计信息 * * * * * * * * * * * * * * * *")
	fmt.Println("客户端发起请求的总数量 = 客户端请求被拦截的数量 + 客户端请求被接受的数量")
	fmt.Println("客户端请求被接受的数量 = 客户端请求下载失败数量 + 客户端请求下载成功数量")
	fmt.Println()
	fmt.Printf("客户端发起请求的总数量: %d个\n", e.scheduler.CalledCount())
	fmt.Printf("客户端请求被拦截的数量: %d个\n", e.scheduler.InterceptCount())
	fmt.Printf("客户端请求被接受的数量: %d个\n", e.scheduler.AcceptedCount())
	fmt.Printf("客户端请求下载失败数量: %d个\n", e.downloader.FailedCount())
	fmt.Printf("客户端请求下载成功数量: %d个\n", e.downloader.CompletedCount())
	fmt.Println("* * * * * * * * * * * * * * * * 统计信息 * * * * * * * * * * * * * * * *")
}
