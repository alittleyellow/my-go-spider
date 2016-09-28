package downloader

import (
	"bytes"
    "common/request"
    "strings"
    "fmt"
    "common/page"
    "github.com/PuerkitoBio/goquery"
    "compress/gzip"
    "golang.org/x/net/html/charset"
    "io"
    "io/ioutil"
    "net/url"
    "net/http"
    "github.com/bitly/go-simplejson"
)

type HttpDownloader struct {

}

func NewHttpDownloader() *HttpDownloader {
	return &HttpDownloader{}
}

func (this *HttpDownloader) Download(req *request.Request) *page.Page {
	var p = page.NewPage(req);
	mtype := req.GetResponceType()
	switch mtype {
	case "html":
		return this.DownloadHtml(p, req)
	case "json":
		return this.DownloadJson(p, req)
	}
	
	return p
}

func (this *HttpDownloader) DownloadHtml(p *page.Page, req *request.Request) *page.Page {

	p, destBody := this.DownloadFile(p, req)
	if !p.IsSuccess() {
		return p
	}
	bodyReader := bytes.NewReader([]byte(destBody))
	doc, err := goquery.NewDocumentFromReader(bodyReader); 
	if err != nil {
		p.SetStatus(true, err.Error())
		return p
	}
	body, err := doc.Html(); 
	if err != nil {
		p.SetStatus(true, err.Error())
        return p
	}

	p.SetBodyStr(body).SetHtmlParser(doc).SetStatus(false, "")

	return p;
}

func (this *HttpDownloader) DownloadJson(p *page.Page, req *request.Request) *page.Page {

	p, destBody := this.DownloadFile(p, req)
	if !p.IsSuccess() {
		return p
	}
	body := []byte(destBody)

	rJson, err := simplejson.NewJson(body)
	if err != nil {
		p.SetStatus(true, err.Error())
		return p
	}


	p.SetBodyStr(string(body)).SetJson(rJson).SetStatus(false, "")

	return p
}

func (this *HttpDownloader) downloadText(p *page.Page, req *request.Request) *page.Page {
    p, destBody := this.DownloadFile(p, req)
    if !p.IsSuccess() {
        return p
    }

    p.SetBodyStr(destBody).SetStatus(false, "")
    return p
}

func (this *HttpDownloader) DownloadFile(p *page.Page, req *request.Request) (*page.Page, string) {
	var err error
	var resp *http.Response

	urlstr := req.GetUrl(); 
	if len(urlstr) == 0 {
		p.SetStatus(true, "url is empty")
		return p, ""
	}
	proxystr := req.GetProxyHost(); 
	if len(proxystr) != 0 {
		resp, err = connectByHttpProxy(p, req)
	} else {
		resp, err = connectByHttp(p, req)
	}
	
	if err != nil {
		return p, ""
	}

	p.SetHeader(resp.Header)
    p.SetCookies(resp.Cookies())

    var bodyStr string

    if resp.Header.Get("Content-Encodind") == "gzip" {
    	bodyStr = this.changeCharsetEncodingAutoGzipSupport(resp.Header.Get("Content-Type"), resp.Body)
    } else {
    	bodyStr = this.changeCharsetEncodingAuto(resp.Header.Get("Content-Type"), resp.Body)
    }
	defer resp.Body.Close()

    return p, bodyStr
}

func (this *HttpDownloader) changeCharsetEncodingAutoGzipSupport(contentTypeStr string, sor io.ReadCloser) string {
	var err error
	gzipReader, err := gzip.NewReader(sor)
	if err != nil {
		return ""
	}
	defer gzipReader.Close()

	destReader, err := charset.NewReader(gzipReader, contentTypeStr)
	if err != nil {
		destReader = sor
	}
	var sorBody []byte
	sorBody, err = ioutil.ReadAll(destReader);
	if err != nil {

	}

	bodystr := string(sorBody)

	return bodystr
}

func (this *HttpDownloader) changeCharsetEncodingAuto(contentTypeStr string, sor io.ReadCloser) string {
	destReader, err := charset.NewReader(sor, contentTypeStr)
	if err != nil {
		destReader = sor
	}
	var sorBody []byte
	sorBody, err = ioutil.ReadAll(destReader); 
	if err != nil {

	}

	bodystr := string (sorBody)

	return bodystr
}

func connectByHttp(p *page.Page, req *request.Request) (*http.Response, error) {
	client := &http.Client{
		CheckRedirect: req.GetRedirectFunc(),
	}
	httpReq, err := http.NewRequest(req.GetMethod(), req.GetUrl(), strings.NewReader(req.GetPostdata()))
	if cookies := req.GetCookies(); cookies != nil {
        for i := range cookies {
        	fmt.Println(cookies[i])
            httpReq.AddCookie(cookies[i])
        }
    }

    resp, err := client.Do(httpReq);
    if err != nil {
    	p.SetStatus(true, err.Error())
    	return nil, err
    }
   
    return resp, nil
}

func connectByHttpProxy(p *page.Page, req *request.Request) (*http.Response, error) {
	proxy, err := url.Parse(req.GetProxyHost())
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	httpReq, _ := http.NewRequest("GET", req.GetUrl(), nil)

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
