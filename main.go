package main

import (
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/yufunny/log-alert/config"
	"github.com/yufunny/log-alert/notify"
	"github.com/yufunny/log-alert/watcher"
	"os"
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
	println(configs.Files)
	for _, file := range configs.Files {
		w := watcher.NewWatcher(file, notifier)
		go w.Watch()
		watchers = append(watchers, w)
	}

	c := cron.New()
	spec := "0 42 21 * * ?"
	c.AddFunc(spec, func() {
		for _, w := range watchers {
			go w.Watch()
		}
	})
	c.Start()

	select {}
}
