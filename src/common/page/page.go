package page

import (
	"github.com/PuerkitoBio/goquery"
    "github.com/bitly/go-simplejson"
    "common/request"
    "net/http"
    "common/page_items"
)

type Page struct {
    isfail   bool
    errormsg string

    // The request is crawled by spider that contains url and relevent information.
    req *request.Request

    // The body is plain text of crawl result.
    body string

    header  http.Header
    cookies []*http.Cookie

    // The docParser is a pointer of goquery boject that contains html result.
    docParser *goquery.Document

    // The jsonMap is the json result.
    jsonMap *simplejson.Json

    // The targetRequests is requests to put into Scheduler.
    targetRequests []*request.Request

    pItems *page_items.PageItems
}

func NewPage(req *request.Request) *Page {
	return &Page{pItems: page_items.NewPageItems(req), req: req}
}

// SetHeader save the header of http responce
func (this *Page) SetHeader(header http.Header) {
    this.header = header
}

// GetHeader returns the header of http responce
func (this *Page) GetHeader() http.Header {
    return this.header
}

// SetHeader save the cookies of http responce
func (this *Page) SetCookies(cookies []*http.Cookie) {
    this.cookies = cookies
}

// GetHeader returns the cookies of http responce
func (this *Page) GetCookies() []*http.Cookie {
    return this.cookies
}

func (this *Page) SetStatus(isfail bool, errormsg string) {
	this.isfail = isfail
	this.errormsg = errormsg
}

func (this *Page) SetBodyStr(body string) *Page {
    this.body = body
    return this
}

func (this *Page) GetBodyStr() string {
    return this.body
}

func (this *Page) SetHtmlParser(doc *goquery.Document) *Page {
    this.docParser = doc
    return this
}

func (this *Page) IsSuccess() bool {
    return !this.isfail
}


func (this *Page) GetRequest() *request.Request {
    return this.req
}

func (this *Page) SetJson(js *simplejson.Json) *Page {
    this.jsonMap = js
    return this
}

func (this *Page) GetJson() *simplejson.Json {
    return this.jsonMap
}











