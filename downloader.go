package gugo

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	MaxRetry         = 5                // 默认最大下载重试次数
	ConnectTimeout   = 10 * time.Second // 默认客户端连接超时时间
	ReadWriteTimeout = 10 * time.Second // 默认客户端读写超时时间
)

var (
	// RetryHTTPCode 默认请求失败时允许重新下载的状态码
	RetryHTTPCode = []int{500, 502, 503, 504, 408}
)

type downloader struct {
	dmu              sync.RWMutex      // 读写锁
	maxRetry         uint32            // 最大下载重试次数
	connectTimeout   time.Duration     // 客户端连接超时时间
	readWriteTimeout time.Duration     // 客户端读写超时时间
	retryHTTPCode    []int             // 下载失败重试状态码
	retryMonitor     map[string]uint32 // 下载失败重试监控器
	*http.Client
	*module
}

func newDownloader() *downloader {
	return &downloader{
		maxRetry:         MaxRetry,
		connectTimeout:   ConnectTimeout,
		readWriteTimeout: ReadWriteTimeout,
		retryHTTPCode:    RetryHTTPCode,
		retryMonitor:     make(map[string]uint32),
		Client:           &http.Client{},
		module:           &module{},
	}
}

// download 统计项：
// 1、下载器正在处理的数量
// 2、客户端请求失败的数量
// 3、客户端请求成功的数量
func (d *downloader) download(req *request, concurrent chan struct{}, reqBuf chan *request, resBuf chan *Response) {
	defer func() { <-concurrent }()
	d.IncrHandlingNumber()
	defer d.DecrHandlingNumber()
	d.Client.Transport = &http.Transport{
		Dial: timeoutDialer(d.connectTimeout, d.readWriteTimeout),
	}
	d.coverClient(req)
	res, err := d.Client.Do(req.Request)
	if err != nil || d.isRetryHTTPCode(res.StatusCode) {
		log.Println(err)
		if d.isNeedRetry(req) {
			go func() { reqBuf <- req }()
			return
		}
		d.IncrFailedCount()
		return
	}
	d.IncrCompletedCount()
	go func() { resBuf <- &Response{res, req} }()
}

// coverClient 覆盖默认客户端
func (d *downloader) coverClient(r *request) {
	if r.meta == nil {
		return
	}
	for _, v := range r.meta {
		if client, ok := v.(*http.Client); ok {
			d.Client = client
			return
		}
	}
}

// isNeedRetry 客户端请求错误是否需要重试
func (d *downloader) isNeedRetry(r *request) bool {
	d.dmu.RLock()
	defer d.dmu.RUnlock()
	var c uint32
	c, ok := d.retryMonitor[r.FingerPrintS()]
	if !ok {
		d.retryMonitor[r.FingerPrintS()] = 1
		return true
	}
	if c >= d.maxRetry-1 {
		return false
	}
	d.retryMonitor[r.FingerPrintS()]++
	return true
}

// isRetryHTTPCode 是否是失败重试状态码
func (d *downloader) isRetryHTTPCode(httpCode int) bool {
	for _, code := range d.retryHTTPCode {
		if code == httpCode {
			return true
		}
	}
	return false
}

// SetMaxRetry 设置最大下载重试次数
func (d *downloader) SetMaxRetry(maxRetry uint32) {
	d.maxRetry = maxRetry
}

// SetRetryHTTPCode 设置下载失败重试状态码
func (d *downloader) SetRetryHTTPCode(httpCode ...int) {
	for _, c := range httpCode {
		if !d.isRetryHTTPCode(c) {
			d.retryHTTPCode = append(d.retryHTTPCode, c)
		}
	}
}

// SetClient 设置客户端
func (d *downloader) SetClient(client *http.Client) {
	d.Client = client
}

// SetConnectTimeout 设置客户端连接超时时间
func (d *downloader) SetConnectTimeout(timeout time.Duration) {
	d.connectTimeout = timeout
}

// SetReadWriteTimeout 设置客户端读写超时时间
func (d *downloader) SetReadWriteTimeout(timeout time.Duration) {
	d.readWriteTimeout = timeout
}

func timeoutDialer(connTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(n, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(n, addr, connTimeout)
		if err != nil {
			return nil, err
		}
		err = conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}
