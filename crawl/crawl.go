package crawl

import (
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
)

type Ext struct {
	*gocrawl.DefaultExtender
}

type ResultSummary struct {
	Http200 int
	Errors int
}

var (
	start time.Time
	responses []*http.Response
	results ResultSummary
)

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	responses = append(responses, res)

	if log.GetLevel() == log.DebugLevel {
		log.WithFields(log.Fields{
			"HeapAlloc": float64(stats.HeapAlloc)/1024.0/1024.0,
			"Sys": float64(stats.Sys)/1024.0/1024.0,
		}).Debug("memory")
	}

	return nil, true
}

func (e *Ext) End(err error) {
	report(responses)
}

func (e *Ext) Error(err *gocrawl.CrawlError) {
	results.Errors++
}

func (e *Ext) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	return true
}

func (e *Ext) ComputeDelay(host string, di *gocrawl.DelayInfo, lastFetch *gocrawl.FetchInfo) time.Duration {
	if lastFetch != nil {
		log.WithFields(log.Fields{
			"url": lastFetch.Ctx.URL(),
			"host": host,
			"duration": lastFetch.Duration,
			"delay-info": di,
			"status": lastFetch.StatusCode,
			"head-request": strconv.FormatBool(lastFetch.IsHeadRequest),
		}).Info("hit")
	}

	return di.OptsDelay
}

func report(r []*http.Response) {
	elapsed := time.Since(start)

	for _, res := range r {
		switch res.StatusCode {
    case 200:
        results.Http200++
    }
	}

	log.WithFields(log.Fields{
		"elapsed": elapsed,
		"total responses": len(r),
		"http-200": results.Http200,
		"errors": results.Errors,
	}).Info("results")
}

func Crawl(url string, external bool, delay int, maxVisits int) {
	start = time.Now()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
	    for sig := range c {
					log.Info(sig)
					report(responses)
					os.Exit(1)
	    }
	}()

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
	crawler := gocrawl.NewCrawlerWithOptions(opts)
	if err := crawler.Run(url); err != nil {
		log.Print(err)
	}
}
