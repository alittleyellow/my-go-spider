package main 
import (
	"common/request"
	"downloader"
	"fmt"
	"io/ioutil"
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
	        Value : "2be13887cdb66026a3509c8c98ec308d",
	        Path : "/",
	        HttpOnly: false,
		},
	}
	req := request.NewRequestWithHeaderFile("http://wx.qima-inc.com/dashboard/user/receivedUserList.json", "html", "", cookie)
	downloaderObject := downloader.NewHttpDownloader();
	resp, _ := downloaderObject.DownloadHtml(req)
	if resp.StatusCode == 200 {
        body, _ := ioutil.ReadAll(resp.Body)
        bodystr := string(body);
        fmt.Println(bodystr)
    }
}