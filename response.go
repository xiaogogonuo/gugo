package gugo

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"unsafe"
)

type Response struct {
	*http.Response
	*request
}

func (r *Response) Valid() bool {
	return r.Response != nil && r.Response.Body != nil
}

func (r *Response) Body() []byte {
	body, _ := io.ReadAll(r.Response.Body)
	r.close()
	return body
}

func (r *Response) Text() string {
	body := r.Body()
	return *(*string)(unsafe.Pointer(&body))
}

func (r *Response) Meta() map[string]interface{} {
	return r.meta
}

func (r *Response) close() {
	_ = r.Response.Body.Close()
}

func (r *Response) XPath(selector string) *goquery.Selection {
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(""))
	return dom.Find(selector)
}
