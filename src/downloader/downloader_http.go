package downloader

import (
    "common/request"
    "net/http"
    "strings"
    "fmt"
)

type HttpDownloader struct {

}

func NewHttpDownloader() *HttpDownloader {
	return &HttpDownloader{}
}

func (this *HttpDownloader) DownloadHtml(req *request.Request) (*http.Response, error) {

	resp, err := this.DownloadFile(req)
	if err != nil {
		fmt.Println(err, 2222)
		return nil, err
	}

	return resp, nil
}

func (this *HttpDownloader) DownloadFile(req *request.Request) (*http.Response, error) {

	resp, err := connectByHttp(req)
	if err != nil {
		fmt.Println(err, 3333)
		return nil, err
	}

	return resp, nil
}

func connectByHttp(req *request.Request) (*http.Response, error) {
	client := &http.Client{
		CheckRedirect: req.GetRedirectFunc(),
	}
	httpReq, err := http.NewRequest(req.GetMethod(), req.GetUrl(), strings.NewReader(req.GetPostdata()))
	if cookies := req.GetCookies(); cookies != nil {
        for i := range cookies {
            httpReq.AddCookie(cookies[i])
        }
    }

    resp, err := client.Do(httpReq)
    if err != nil {
    	fmt.Println(err, 4444)
    	return nil, err	
    }

    return resp, nil
}