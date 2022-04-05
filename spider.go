package gugo

const (
	ResponseBufCap     = 1 << 12 // 默认响应队列容量
	ConcurrentResponse = 1 << 10 // 默认响应处理的并发量
)

type spider struct {
	resBuf             chan *Response // 响应队列
	concurrentResponse chan struct{}  // 响应并发控制
	*module
}

func newSpider() *spider {
	return &spider{
		resBuf:             make(chan *Response, ResponseBufCap),
		concurrentResponse: make(chan struct{}, ConcurrentResponse),
		module:             &module{},
	}
}

// parse 爬虫解析
func (s *spider) response(res *Response) {
	s.IncrHandlingNumber()
	defer s.DecrHandlingNumber()
	defer func() { <-s.concurrentResponse }()
	res.parser(res)
}

// SetResponseBufCap 设置响应队列容量
func (s *spider) SetResponseBufCap(n uint32) {
	s.resBuf = make(chan *Response, n)
}

// SetConcurrentResponse 设置响应处理的并发量
func (s *spider) SetConcurrentResponse(n uint32) {
	s.concurrentResponse = make(chan struct{}, n)
}
