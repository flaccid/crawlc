package main

import (
	"os"

	"github.com/flaccid/crawlc/crawl"
	"github.com/urfave/cli"
	log "github.com/Sirupsen/logrus"
)

var (
	VERSION = "v0.0.0-dev"
)

func beforeApp(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	if len(c.Args().Get(0)) < 1 {
		log.Fatal("please provide an url to crawl")
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "crawlc"
	app.Version = VERSION
	app.Usage = "a website crawler, frontend to gocrawl"
	app.Action = start
	app.Before = beforeApp
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "external,e",
			Usage: "crawl external urls",
		},
		cli.IntFlag{
			Name:  "crawl-delay,d",
			Usage: "crawl delay in seconds",
			Value: 1,
		},
		cli.IntFlag{
			Name:  "max-visits,m",
			Usage: "maximum number of visits",
			Value: 100,
		},
		cli.BoolFlag{
			Name:  "debug,D",
			Usage: "run in debug mode",
		},
	}
	app.Run(os.Args)
}

func start(c *cli.Context) error {
	crawl.Crawl(c.Args().Get(0), c.Bool("external"), c.Int("crawl-delay"), c.Int("max-visits"))

	return nil
}
