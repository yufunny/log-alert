package watcher

import (
	"fmt"
	"github.com/papertrail/go-tail/follower"
	"github.com/sirupsen/logrus"
	"github.com/yufunny/log-alert/notify"
	"io"
	"os"
	"regexp"
	"time"
	"github.com/yufunny/log-alert/config"
	"strings"
)

type Watcher struct {
	file     string
	desc     string
	rules	[]rule
	analyser string

	notifier notify.Notify
	handler * follower.Follower
	live uint64
}

type rule struct {
	duration time.Duration
	times    int
	interval time.Duration
	count    int
	sent     bool
	text     []string
}

func NewWatcher(fileConfig config.FileConfig, notifier notify.Notify) *Watcher {
	watcher :=  &Watcher{
		file: fileConfig.File,
		desc: fileConfig.Desc,
		analyser: fileConfig.Analyser,
		notifier: notifier,
	}
	return watcher
}

func parseFile(raw string) string {
	if strings.Index(raw, "%Y") > -1 {
		logrus.Debugf("time format %d", time.Now().Year())
		//strings.Replace(raw, "%Y", time.Now().Year())
	}
	return raw
}

func parseRule(ruleConfig config.RuleConfig) *rule {

	duration, _ := time.ParseDuration(ruleConfig.Duration)
	interval, _ := time.ParseDuration(ruleConfig.Interval)
  return &rule{
  	duration: duration,
  	interval: interval,
  	times: ruleConfig.Times,
  	count: 0,
  	sent: false,
  	text: make([]string, 0),
  }
}

func (w *Watcher) Watch() {
	parsedFile := parseFile(w.file)
	w.handler, _ = follower.New(parsedFile, follower.Config{
		Whence: io.SeekEnd,
		Offset: 0,
		Reopen: true,
	})

	go w.tick()
	//reg := regexp.MustCompile(w.Rule)
	//
	//for line := range w.handler.Lines() {
	//	if reg.Match(line.Bytes()) {
	//		logrus.Debugf("%s:%s", w.desc, string(line.Bytes()))
	//		w.Text = append(w.Text, string(line.Bytes()))
	//		if len(w.Text) >= w.Times && !w.Sent {
	//			w.Notifier.Send(w.Desc, w.Text...)
	//			if w.Interval.Nanoseconds() > 0 {
	//				w.Sent = true
	//			}
	//			w.Text = make([]string, 0)
	//		}
	//	}
	//}
	//
	//if t.Err() != nil {
	//	fmt.Fprintln(os.Stderr, t.Err())
	//}
}

func (w *Watcher) tick() {
	clock := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-clock.C:
				w.live++
				if w.live == ^uint64(0) {
					w.live = 0
				}
				//for w.Rules
			}
		}
	}()
}
