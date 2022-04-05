package gugo

import (
	"github.com/bits-and-blooms/bloom/v3"
	"sync"
	"time"
)

const (
	FalsePositive     = 0.01    // 默认过滤错误容忍率
	EstimateRequest   = 1000000 // 默认估计100万请求
	RequestBufferCap  = 1 << 12 // 默认请求队列容量
	ConcurrentRequest = 1 << 10 // 默认并发请求数
)

type scheduler struct {
	smu               sync.Mutex
	filter            *bloom.BloomFilter  // 布隆过滤器
	domain            map[string]struct{} // 可用域名
	reqBuf            chan *request       // 请求队列
	concurrentRequest chan struct{}       // 请求并发控制
	duration          time.Duration       // 下载延时
	*module
}

func newScheduler() *scheduler {
	return &scheduler{
		filter:            bloom.NewWithEstimates(EstimateRequest, FalsePositive),
		domain:            make(map[string]struct{}),
		reqBuf:            make(chan *request, RequestBufferCap),
		concurrentRequest: make(chan struct{}, ConcurrentRequest),
		module:            &module{},
	}
}

// ask 统计项：
// 1、调度器正在处理的数量
// 2、客户端发起请求的数量
// 3、请求被拦截过滤的数量
// 4、请求被接受下载的数量
func (s *scheduler) ask(r *request) {
	s.IncrHandlingNumber()
	defer s.DecrHandlingNumber()
	s.IncrCalledCount()
	if !s.isAcceptedRequest(r) {
		s.IncrInterceptCount()
		return
	}
	s.IncrAcceptedCount()
	time.Sleep(s.duration)
	go func() { s.reqBuf <- r }()
}

// isAcceptedRequest 判断请求是否可访问
func (s *scheduler) isAcceptedRequest(r *request) bool {
	return r.Valid() &&
		s.isAcceptedDomain(r) &&
		s.isAcceptedSchema(r) &&
		s.isUniqueRequest(r)
}

// isUniqueRequest 判断请求是否重复
func (s *scheduler) isUniqueRequest(r *request) bool {
	s.smu.Lock()
	defer s.smu.Unlock()
	if !s.filter.Test(r.FingerPrint()) {
		s.filter.Add(r.FingerPrint())
		return true
	}
	return false
}

// isAcceptedSchema 判断请求协议是否可访问
func (s *scheduler) isAcceptedSchema(r *request) bool {
	if r.Schema() == "http" || r.Schema() == "https" {
		return true
	}
	return false
}

// isAcceptedDomain 判断请求域名是否可访问
func (s *scheduler) isAcceptedDomain(r *request) bool {
	if len(s.domain) == 0 {
		return true
	}
	if _, ok := s.domain[r.Host()]; ok {
		return true
	}
	return false
}

// SetDomain 设置可访问的域名
func (s *scheduler) SetDomain(domain ...string) {
	for _, v := range domain {
		if _, ok := s.domain[v]; !ok {
			s.domain[v] = struct{}{}
		}
	}
}

// SetRequestBufCap 设置请求队列容量
func (s *scheduler) SetRequestBufCap(n uint32) {
	s.reqBuf = make(chan *request, n)
}

// SetConcurrentRequest 设置请求处理的并发量
func (s *scheduler) SetConcurrentRequest(n uint32) {
	s.concurrentRequest = make(chan struct{}, n)
}
