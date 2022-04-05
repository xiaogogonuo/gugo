package gugo

const (
	PipelineBufCap     = 1 << 12 // 默认数据队列容量
	ConcurrentPipeline = 1 << 10 // 默认数据处理的并发量
)

type pipeline struct {
	pipeBuf            chan interface{} // 数据队列
	concurrentPipeline chan struct{}    // 数据并发控制
}

func newPipeline() *pipeline {
	return &pipeline{
		pipeBuf:            make(chan interface{}, PipelineBufCap),
		concurrentPipeline: make(chan struct{}, ConcurrentPipeline),
	}
}

// push 接受客户端推送的数据
func (p *pipeline) push(item interface{}) {
	go func() { p.pipeBuf <- item }()
}

// pull 客户端拉取数据
func (p *pipeline) pull() chan interface{} {
	return p.pipeBuf
}

// Empty 数据管道是否空
func (p *pipeline) Empty() bool {
	return len(p.pipeBuf) == 0
}

// SetPipelineBufCap 设置数据队列容量
func (p *pipeline) SetPipelineBufCap(n uint32) {
	p.pipeBuf = make(chan interface{}, n)
}

// SetConcurrentPipeline 设置数据处理的并发量
func (p *pipeline) SetConcurrentPipeline(n uint32) {
	p.concurrentPipeline = make(chan struct{}, n)
}
