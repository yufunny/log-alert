package watcher

import (
	"github.com/papertrail/go-tail/follower"
	"github.com/yufunny/log-alert/notify"
	"io"
	"regexp"
	"time"
	"github.com/yufunny/log-alert/config"
	"strings"
	"strconv"
	"fmt"
	"os"
	"github.com/sirupsen/logrus"
)

type Watcher struct {
	file     string
	desc     string
	rules	[]*rule
	boundRegexp *regexp.Regexp

	notifier notify.Notify
	handler * follower.Follower
	live uint64
	piece []string
	clock *time.Ticker
}

type rule struct {
	ruleRegexp *regexp.Regexp
	desc     string
	duration uint64
	times    int
	interval uint64
	count    int
	sent     bool
	text     []string
	receivers []string
}

func NewWatcher(fileConfig config.FileConfig, notifier notify.Notify) *Watcher {
	rules := make([]*rule, 0)
	for _, rule := range fileConfig.Rules {
		parsedRule := parseRule(rule)
		rules = append(rules, parsedRule)
	}
	var boundRegexp *regexp.Regexp
	if fileConfig.Bound != "" {
		boundRegexp = regexp.MustCompile(fileConfig.Bound)
	} else {
		boundRegexp = nil
	}
	watcher :=  &Watcher{
		file: fileConfig.File,
		desc: fileConfig.Desc,
		boundRegexp: boundRegexp,
		notifier: notifier,
		rules: rules,
	}
	go watcher.tick()
	return watcher
}

func parseFile(raw string) string {
	if strings.Index(raw, "%Y") > -1 {
		raw = strings.Replace(raw, "%Y", strconv.Itoa(time.Now().Year()), 1)
	}
	if strings.Index(raw, "%m") > -1 {
		month := fmt.Sprintf("%02d", time.Now().Month())
		raw = strings.Replace(raw, "%m", month, 1)
	}
	if strings.Index(raw, "%d") > -1 {
		day := fmt.Sprintf("%02d", time.Now().Day())
		raw = strings.Replace(raw, "%d", day, 1)
	}
	return raw
}

func parseRule(ruleConfig config.RuleConfig) *rule {
	duration, _ := time.ParseDuration(ruleConfig.Duration)
	interval, _ := time.ParseDuration(ruleConfig.Interval)
	  return &rule{
		ruleRegexp: regexp.MustCompile(ruleConfig.Rule),
		desc: ruleConfig.Desc,
		duration: uint64(duration.Seconds()),
		interval: uint64(interval.Seconds()),
		times: ruleConfig.Times,
		count: 0,
		sent: false,
		text: make([]string, 0),
		receivers: ruleConfig.Receiver,
	  }
}

func (w *Watcher) Watch() {
	if w.handler != nil {
		w.handler.Close()
	}
	parsedFile := parseFile(w.file)
	for {
		_, err := os.Stat(parsedFile)
		if err == nil {
			break
		}
		logrus.Infof("文件:%s 不存在", parsedFile)
		time.Sleep(time.Minute)
	}
	w.handler, _ = follower.New(parsedFile, follower.Config{
		Whence: io.SeekEnd,
		Offset: 0,
		Reopen: true,
	})

	logrus.Infof("start listening: %s", parsedFile)

	for line := range w.handler.Lines() {
		piece := w.parsePiece(line.String())
		if len(piece) == 0 {
			continue
		}

		for _, rule := range w.rules {
			if rule.ruleRegexp.Match([]byte(piece[0])) {
				if !rule.sent {
					rule.text = append(rule.text, piece...)
					rule.count++
					if rule.count >= rule.times {
						w.notifier.Send(rule.receivers, "[" + w.desc +"]" + rule.desc, rule.text...)
						if rule.interval > 0 {
							rule.sent = true
						}
						rule.text = make([]string, 0)
						rule.count = 0
					}
				}
			}
		}

	}
}

func (w *Watcher) parsePiece(line string) []string {
	if w.boundRegexp == nil {
		return []string{line}
	}
	if ! w.boundRegexp.Match([]byte(line)) {
		if len(w.piece) != 0 {
			w.piece = append(w.piece, line)
		}
		return []string{}
	} else {
		ret := w.piece
		w.piece = []string{line}
		return ret
	}
}

func (w *Watcher) tick() {
	w.clock = time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-w.clock.C:
				w.live++
				if w.live == ^uint64(0) {
					w.live = 0
				}
				for _, rule := range w.rules {
					if rule.interval> 0 && w.live % rule.interval == 0 {
						logrus.Debugf("internal clear")
						rule.sent = false
					}
					if rule.duration> 0 && w.live % rule.duration == 0 {
						logrus.Debugf("duration clear")
						rule.count = 0
						rule.text = make([]string, 0)
					}
				}
			}
		}
	}()
}
