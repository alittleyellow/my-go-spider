package main
import (
	"spider"
	"common/page"
	"common/request"
	"fmt"
	"strings"
	"net/http"
)
type MyPageProcesser struct {

}

func NewMyPageProcesser() *MyPageProcesser{
	return &MyPageProcesser{}
}

func (this *MyPageProcesser) Process(p *page.Page) {
	if !p.IsSuccess() {
		return 
	}
	query := p.GetHtmlParser()
	name := query.Find("h1").Text()
    name = strings.Trim(name, " \t\n")
    fmt.Println(name)
}

func (this *MyPageProcesser) Finish() {
	fmt.Printf("TODO:before end spider \r\n")
}

func main() {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		sp := spider.NewSpider(NewMyPageProcesser(), "myFirstTaskName")
	 	//多个请求
	 	urls := []string{
	        "http://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	        "http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
	    }
	 	var reqs []*request.Request
	    for _, url := range urls {
	        req := request.NewRequest(url, "html", "GET", "", "", nil, nil, "", nil, nil)
	        reqs = append(reqs, req)
	    }
	    sp.SetThreadnum(5).GetAllByRequest(reqs)
	})

	http.ListenAndServe(":3002", nil)
	
}











