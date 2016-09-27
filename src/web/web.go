package main 
import (
	"spider"
	"common/page"
	"common/request"
	"fmt"
	"strings"
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
    summary := query.Find(".card-summary-content .para").Text()
    summary = strings.Trim(summary, " \t\n")

    // the entity we want to save by Pipeline
    p.AddField("name", name)
    p.AddField("summary", summary)
}

func (this *MyPageProcesser) Finish() {
	fmt.Printf("TODO:before end spider \r\n")
}


func main() {
	sp := spider.NewSpider(NewMyPageProcesser(), "myFirstTaskName")
	req := request.NewRequest("https://bbs.youzan.com/forum.php?mod=viewthread&tid=543951", "html", "GET", "", "", nil, nil, "", nil, nil)
	pageItems := sp.GetByRequest(req)
	for name, _ := range pageItems.GetAll() {
		fmt.Println(name)
    }
}











