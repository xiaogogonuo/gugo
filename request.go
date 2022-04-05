package gugo

import (
	"github.com/xiaogogonuo/gugo/pkg/crypto"
	"io"
	"net/http"
	"unsafe"
)

type Parser func(*Response)

type request struct {
	*http.Request
	parser Parser
	meta   map[string]interface{}
}

func (r *request) Valid() bool {
	return r.Request != nil && r.Request.URL != nil && r.parser != nil
}

func (r *request) URL() string {
	return r.Request.URL.String()
}

func (r *request) Host() string {
	return r.Request.Host
}

func (r *request) Body() []byte {
	if r.Request.Body == nil {
		return []byte{}
	}
	body, _ := io.ReadAll(r.Request.Body)
	return body
}

func (r *request) Method() string {
	return r.Request.Method
}

func (r *request) Schema() string {
	return r.Request.URL.Scheme
}

// FingerPrint 请求指纹：sha1(请求体+请求URL+请求方法)
func (r *request) FingerPrint() []byte {
	finger := r.Body()
	url, method := r.URL(), r.Method()
	finger = append(finger, *(*[]byte)(unsafe.Pointer(&url))...)
	finger = append(finger, *(*[]byte)(unsafe.Pointer(&method))...)
	return crypto.SHA1Encrypt2Byte(&finger)
}

// FingerPrintS 请求指纹字符串
func (r *request) FingerPrintS() string {
	finger := r.FingerPrint()
	return *(*string)(unsafe.Pointer(&finger))
}
