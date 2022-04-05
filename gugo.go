package gugo

import (
	"net/http"
)

type GuGo struct {
	*engine
}

// CreateGuGo 😄创建谷歌😄
func CreateGuGo() *GuGo {
	return &GuGo{engine: newEngine()}
}

// Request 简易版GET请求
func (g *GuGo) Request(url string, parser Parser, meta map[string]interface{}) {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	g.NativeRequest(request, parser, meta)
}

// NativeRequest 原生请求，客户端自定义
func (g *GuGo) NativeRequest(r *http.Request, parser Parser, meta map[string]interface{}) {
	g.ask(&request{r, parser, meta})
}

// Push 客户端发送数据
func (g *GuGo) Push(item interface{}) {
	g.push(item)
}

// Pull 客户端下载数据
func (g *GuGo) Pull() chan interface{} {
	return g.pull()
}

// GooGol 谷歌运行入口
func (g *GuGo) GooGol() {
	g.coordinate()
}
