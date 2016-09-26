package main 
import (
	"common/request"
	"downloader"
	"fmt"
	// "io/ioutil"
	"net/http"
)

func main() {
	cookie := []*http.Cookie{
		&http.Cookie{
			Name : "_servant_value",
	        Value : "149",
	        Path : "/",
	        HttpOnly: false,
		},
		&http.Cookie{
			Name : "_servant_key",
	        Value : "41330d4ffd845bc32ac0942852706a1c",
	        Path : "/",
	        HttpOnly: false,
		},
	}
	req := request.NewRequestWithHeaderFile("http://wx.qima-inc.com/dashboard/user/receivedUserList.json", "json", "", cookie)
	downloaderObject := downloader.NewHttpDownloader();
	p := downloaderObject.Download(req)
	fmt.Println(p.GetJson());
}