package crawl

import (
	"net/http"
	"runtime"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
)

type Ext struct {
	*gocrawl.DefaultExtender
}

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	log.Infof("Visit: %s\n", ctx.URL())

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	log.Debugf("HeapAlloc=%02fMB; Sys=%02fMB\n", float64(stats.HeapAlloc)/1024.0/1024.0, float64(stats.Sys)/1024.0/1024.0)

	return nil, true
}

func (e *Ext) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	return true
}

func Crawl(url string, external bool, delay int, maxVisits int) {
	ext := &Ext{&gocrawl.DefaultExtender{}}
	opts := gocrawl.NewOptions(ext)

	if external {
		opts.SameHostOnly = false
	} else {
		opts.SameHostOnly = true
	}

	opts.CrawlDelay = time.Duration(delay) * time.Second
	opts.LogFlags = gocrawl.LogError
	opts.MaxVisits = maxVisits

	log.Debug(opts, url)

	log.Print("starting crawl...")
	c := gocrawl.NewCrawlerWithOptions(opts)
	if err := c.Run(url); err != nil {
		log.Print(err)
	}
}
