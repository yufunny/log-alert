package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/yufunny/log-alert/config"
	"github.com/yufunny/log-alert/notify"
	"github.com/yufunny/log-alert/watcher"
	"os"
	"time"
)

var (
	configPath string
	configs    *config.SystemConfig
)

func main() {
	log.SetLevel(log.DebugLevel)
	{
		app := cli.NewApp()
		app.Name = "log-alert"
		app.Usage = "alert by log file"
		app.Version = "0.1.0"
		app.Authors = []cli.Author{
			{
				Name:  "yufu",
				Email: "mxy@yufu.fun",
			},
		}
		app.Flags = []cli.Flag{
			cli.StringFlag{
				Name:        "config, c",
				Usage:       "config path",
				Value:       "config.yaml",
				Destination: &configPath,
			},
		}
		app.Action = run

		log.Infof("[MAIN]Run start")
		err := app.Run(os.Args)
		if nil != err {
			log.Errorf("[MAIN]Run error; " + err.Error())
		}
	}
}

func run(_ *cli.Context) {
	var err error
	configs, err = config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("[main]parse config error:%s", err.Error())
	}
	if configs.Mode == "release" {
		log.SetLevel(log.InfoLevel)
	}

	notifier := &notify.MailNotify{
		Url:       configs.Notify.Url,
		Receivers: configs.Receiver,
	}
	watchers := make([]*watcher.Watcher, 0)
	for _, rule := range configs.Rules {
		duration, _ := time.ParseDuration(rule.Duration)
		interval, _ := time.ParseDuration(rule.Interval)
		watch := &watcher.Watcher{
			File:     rule.File,
			Rule:     rule.Rule,
			Desc:     rule.Desc,
			Duration: duration,
			Times:    rule.Times,
			Interval: interval,
			Notifier: notifier,
			Count:    0,
			Sent:     false,
		}
		go watch.Watch()
		watchers = append(watchers, watch)
	}

	for {
		time.Sleep(time.Second)
	}
}
