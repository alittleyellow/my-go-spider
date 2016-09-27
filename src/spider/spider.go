package spider

import (
	"page_processer"
	"downloader"
	"scheduler"
	"pipeline"
	"common/resource_manage"
	"common/request"
	"common/page_items"
	"common/page"
	"time"
	"math/rand"
	"fmt"
)

type Spider struct {
	taskname string

    pPageProcesser page_processer.PageProcesser

    pDownloader *downloader.HttpDownloader

    pScheduler scheduler.Scheduler

    pPiplelines []pipeline.Pipeline

    mc resource_manage.ResourceManage

    threadnum uint

    exitWhenComplete bool

    // Sleeptype can be fixed or rand.
    startSleeptime uint
    endSleeptime   uint
    sleeptype      string
}

func NewSpider(pageinst page_processer.PageProcesser, taskname string) *Spider {
	fmt.Println("new spider")
	ap := &Spider{pPageProcesser: pageinst, taskname: taskname}
	ap.exitWhenComplete = true
	ap.sleeptype = "fixed"
	ap.startSleeptime = 0

	if ap.pScheduler == nil {
		ap.SetScheduler(scheduler.NewQueueScheduler(false))
	}

	if ap.pDownloader == nil {
		ap.SetDownloader(downloader.NewHttpDownloader())
	}

	ap.pPiplelines = make([]pipeline.Pipeline, 0)

	return ap
}

func (this *Spider) SetScheduler(s scheduler.Scheduler) *Spider {
	this.pScheduler = s
    return this
} 

func (this *Spider) SetDownloader(d *downloader.HttpDownloader) *Spider {
	this.pDownloader = d
	return this
}

func (this *Spider) Taskname() string {
    return this.taskname
}

func (this *Spider) Get(url string, respType string) *page_items.PageItems {
    req := request.NewRequest(url, respType, "GET", "", "", nil, nil, "", nil, nil)
    return this.GetByRequest(req)
}

func (this *Spider) GetByRequest(req *request.Request) *page_items.PageItems {
	var reqs []*request.Request
    reqs = append(reqs, req)
    items := this.GetAllByRequest(reqs)

    if len(items) != 0 {
        return items[0]
    }

    return nil
}

func (this *Spider) GetAllByRequest(reqs []*request.Request) []*page_items.PageItems {
	for _, req := range reqs {
		this.AddRequest(req)
	}
	pip := pipeline.NewCollectPipelinePageItems()
	this.AddPipeline(pip)

	this.Run()

	return pip.GetCollected()
}

func (this *Spider) AddRequest(req *request.Request) *Spider {
	if req == nil {
		return this
	} else if req.GetUrl() == "" {
		return this
	}

	this.pScheduler.Push(req)
	return this
}

func (this *Spider) AddPipeline(pip pipeline.Pipeline) *Spider {
	this.pPiplelines = append(this.pPiplelines, pip)
	return this
}

func (this *Spider) Run() {
	if this.threadnum == 0 {
		this.threadnum = 1
	}	

	this.mc = resource_manage.NewResourceManageChan(this.threadnum)

	for {
		req := this.pScheduler.Poll()
		if this.mc.Has() == 0 && req == nil && this.exitWhenComplete {
			this.pPageProcesser.Finish()
			break
		} else if req == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		this.mc.GetOne()

		go func(req *request.Request) {
			defer this.mc.FreeOne()
			this.pageProcess(req)
		}(req)
	}

	this.close()
}

func (this *Spider) pageProcess(req *request.Request) {
	var p *page.Page 

	//download page
	for i := 0; i < 3; i++ { //如果失败 重试三次
		this.sleep()
		p = this.pDownloader.Download(req)
		if p.IsSuccess() {
			break
		}
	}

	if !p.IsSuccess() {
		return 
	}

	this.pPageProcesser.Process(p)
	for _, req := range p.GetTargetRequests() {
		this.AddRequest(req)
	}

	if !p.GetSkip() {
		for _, pip := range this.pPiplelines {
			pip.Process(p.GetPageItems(), this)
		}
	}

}

func (this *Spider) close() {
	this.SetScheduler(scheduler.NewQueueScheduler(false))
    this.SetDownloader(downloader.NewHttpDownloader())
    this.pPiplelines = make([]pipeline.Pipeline, 0)
    this.exitWhenComplete = true
}	

func (this *Spider) sleep() {
	if this.sleeptype == "fixed" {
		time.Sleep(time.Duration(this.startSleeptime) * time.Millisecond)
	} else if this.sleeptype == "rand" {
		sleeptime := rand.Intn(int(this.endSleeptime-this.startSleeptime)) + int(this.startSleeptime)
        time.Sleep(time.Duration(sleeptime) * time.Millisecond)
	}
}

func (this *Spider) AddUrl(url string, respType string) *Spider {
    req := request.NewRequest(url, respType, "GET", "", "", nil, nil, "", nil, nil)
    this.AddRequest(req)
    return this
}

func (this *Spider) AddUrlEx(url string, respType string, headerFile string, proxyHost string) *Spider {
    req := request.NewRequest(url, respType, "GET", "", "", nil, nil, "", nil, nil)
    this.AddRequest(req.AddHeaderFile(headerFile).AddProxyHost(proxyHost))
    return this
}

func (this *Spider) AddUrlWithHeaderFile(url string, respType string, headerFile string) *Spider {
    req := request.NewRequestWithHeaderFile(url, respType, headerFile, nil)
    this.AddRequest(req)
    return this
}

func (this *Spider) AddUrls(urls []string, respType string) *Spider {
    for _, url := range urls {
        req := request.NewRequest(url, respType, "GET", "", "", nil, nil, "", nil, nil)
        this.AddRequest(req)
    }
    return this
}

func (this *Spider) AddUrlsWithHeaderFile(urls []string, respType string, headerFile string) *Spider {
    for _, url := range urls {
        req := request.NewRequestWithHeaderFile(url, respType, headerFile, nil)
        this.AddRequest(req)
    }
    return this
}

func (this *Spider) AddUrlsEx(urls []string, respType string, headerFile string, proxyHost string) *Spider {
    for _, url := range urls {
        req := request.NewRequest(url, respType, "GET", "", "", nil, nil, "", nil, nil)
        this.AddRequest(req.AddHeaderFile(headerFile).AddProxyHost(proxyHost))
    }
    return this
}