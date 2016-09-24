package util 

import (
	"strings"
	"fmt"
	"regexp"
	"os"
)

func JsonpToJson(json string) string {
	fmt.Println(json)
	start := strings.Index(json, "{")
	end := strings.LastIndex(json, "}")
	startP := strings.Index(json, "]")
	if startP > 0 && start > startP {
		start = startP
		end = strings.LastIndex(json, "]")
	}
	if end > start && end != -1 && start != -1 {
		json = json[start : end+1]
	}
	json = strings.Replace(json, "\\'", "", -1)
	regDetail, _ := regexp.Compile("([^\\s\\:{\\,\\d\"]+|[a-z][a-z\\d]*)\\s*\\:")
	jsonp := regDetail.ReplaceAllString(json, "\"$1\":")

	return jsonp;
}

func GetWDPath() string {
	wd := os.Getenv("GOPATH")
	if wd == "" {
		panic("GOPATH IS NOT SET IN ENV")
	}

	return wd
}