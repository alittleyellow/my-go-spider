package request

import (
	"net/http"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"os"
)

type Request struct {
	//请求的url
	Url string
	//响应格式 json xml jsonp... 
	RespType string
	// GET POST
    Method string
	//post data
	Postdata string
	//url打标区分
	Urltag string
	//http header
	Header http.Header
	//
	Cookies []*http.Cookie
	//proxy host example='localhost:80'
	ProxyHost string
	//
	checkRedirect func(req *http.Request, via []*http.Request) error
	//
	Meta interface{}
}

func NewRequest(url, respType, method, postData, urlTag string, header http.Header, cookies []*http.Cookie, proxyHost string, checkRedirect func(req *http.Request, via []*http.Request) error, meta interface{}) *Request {
	
	return &Request{url, respType, method, postData, urlTag, header, cookies, proxyHost, checkRedirect, meta}
}

func NewRequestWithHeaderFile(url, respType, headerFile string, cookies []*http.Cookie) *Request {
	_, err := os.Stat(headerFile)
	if err != nil {
		return &Request{url, respType, "GET", "", "", nil, cookies, "", nil, nil}
	}
	h := ReadHeaderFromFile(headerFile)

	return &Request{url, respType, "GET", "", "", h, nil, "", nil, nil}
}

func ReadHeaderFromFile(headerFile string) http.Header {
	b, err := ioutil.ReadFile(headerFile)
	if err != nil {
		return nil
	}
	js, _ := simplejson.NewJson(b)
	header := make(http.Header)
	header.Add("User-Agent", js.Get("User-Agent").MustString())
	header.Add("Referer", js.Get("Referer").MustString())
    header.Add("Cookie", js.Get("Cookie").MustString())
    header.Add("Cache-Control", "max-age=0")
    header.Add("Connection", "keep-alive")

    return header
}

func (this *Request) AddHeaderFile(headerFile string) *Request {
    _, err := os.Stat(headerFile)
    if err != nil {
        return this
    }
    h := ReadHeaderFromFile(headerFile)
    this.Header = h

    return this
}

// @host  http://localhost:8765/
func (this *Request) AddProxyHost(host string) *Request {
    this.ProxyHost = host

    return this
}

func (this *Request) GetUrl() string {

    return this.Url
}

func (this *Request) GetUrlTag() string {

    return this.Urltag
}

func (this *Request) GetMethod() string {

    return this.Method
}

func (this *Request) GetPostdata() string {

    return this.Postdata
}

func (this *Request) GetHeader() http.Header {

    return this.Header
}

func (this *Request) GetCookies() []*http.Cookie {

    return this.Cookies
}

func (this *Request) GetProxyHost() string {

    return this.ProxyHost
}

func (this *Request) GetResponceType() string {

    return this.RespType
}

func (this *Request) GetRedirectFunc() func(req *http.Request, via []*http.Request) error {

    return this.checkRedirect
}

func (this *Request) GetMeta() interface{} {

    return this.Meta
}
