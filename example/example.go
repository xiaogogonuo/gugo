package main

import (
	"fmt"
	"github.com/xiaogogonuo/gugo"
	"net/http"
)

// Film 客户端自定义模型
type Film struct {
	score uint
	title string
}

type MySpider struct {
	*gugo.GuGo
	block chan struct{} // 客户端数据处理阻塞控制器，防止爬虫结束，数据处理被迫停止
}

func (ms *MySpider) Parse1(response *gugo.Response) {
	// 获取响应数据的字节流形式
	_ = response.Body()
	// 获取响应数据的字符串形式
	_ = response.Text()
	// 获取请求携带过来的元数据
	_ = response.Meta()
	// 获取响应所对应的请求链接
	_ = response.URL()
	// 获取响应所对应的请求方法
	_ = response.Method()

	// 模拟发送从页面提取的数据
	ms.Push(Film{score: 5, title: "肖申克的救赎"})

	// 模拟发送从页面提取的新链接
	// Request方法是简易版的GET请求，客户端只需传入URL，自定义解析函数，元数据即可
	// 新链接：必填
	// 解析器：必填
	// 元数据：可选
	ms.Request("https://www.douban.com", ms.Parse2, map[string]interface{}{"parser": "Parser1"})
}

func (ms *MySpider) Parse2(response *gugo.Response) {
	// 获取从Parser1传入的元数据
	_ = response.Meta()

	// 模拟发送从页面提取的数据
	ms.Push(Film{score: 4, title: "越狱"})

	// 模拟发送从页面提取的新链接
	// NativeRequest方法是原生版的http请求，客户端需要传入自定义*http.Request、自定义解析函数、元数据
	// 新请求：必填
	// 解析器：必填
	// 元数据：可选
	customGetRequest, _ := http.NewRequest(http.MethodGet, "https://www.tencent.com", nil)
	ms.NativeRequest(customGetRequest, ms.Parser3, map[string]interface{}{"parser": "Parser2"})

	customPostRequest, _ := http.NewRequest(http.MethodPost, "https://www.abc.com", nil)
	ms.NativeRequest(customPostRequest, ms.Parser3, map[string]interface{}{"parser": "Parser2"})
}

func (ms *MySpider) Parser3(response *gugo.Response) {
	// 获取从Parser2传入的元数据
	_ = response.Meta()

	// 模拟发送从页面提取的数据
	ms.Push(Film{score: 3, title: "复仇者联盟"})
}

// ProcessItem 客户端自定义数据处理器
func (ms *MySpider) ProcessItem() {
	for item := range ms.Pull() {
		fmt.Println(item)
	}
	ms.block <- struct{}{}
}

func main() {
	// 1、创建我的爬虫
	ms := &MySpider{
		GuGo:  gugo.CreateGuGo(),
		block: make(chan struct{}),
	}
	// 2、发送初始请求
	ms.Request("https://www.baidu.com", ms.Parse1, nil)
	// 3、异步处理客户端数据(保存到文件、数据库等操作)
	go ms.ProcessItem()
	// 4、启动并发爬虫系统
	ms.GooGol()
	// 5、客户端数据处理阻塞
	<-ms.block
}
