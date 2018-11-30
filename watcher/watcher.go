package watcher

import (
	"github.com/papertrail/go-tail/follower"
	"github.com/sirupsen/logrus"
	"github.com/yufunny/log-alert/notify"
	"io"
	"regexp"
	"time"
	"github.com/yufunny/log-alert/config"
	"strings"
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
		ruleRegexp: regexp.MustCompile(ruleConfig.Rule),
		desc: ruleConfig.Desc,
		duration: uint64(duration.Seconds()),
		interval: uint64(interval.Seconds()),
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
	for line := range w.handler.Lines() {
		logrus.Debugf("%s:%s", w.desc, string(line.Bytes()))
		piece := w.parsePiece(line.String())
		if len(piece) == 0 {
			continue
		}

		for _, rule := range w.rules {
			if rule.ruleRegexp.Match(line.Bytes()) {
				rule.text = append(rule.text, string(line.Bytes()))
				if len(rule.text) >= rule.times && !rule.sent {
					w.notifier.Send(rule.desc, rule.text...)
					if rule.interval > 0 {
						rule.sent = true
					}
					rule.text = make([]string, 0)
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
	clock := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-clock.C:
				w.live++
				if w.live == ^uint64(0) {
					w.live = 0
				}
				for _, rule := range w.rules {
					if rule.interval> 0 && w.live % rule.interval == 0 {
						rule.sent = false
					}
					if rule.duration> 0 && w.live % rule.duration == 0 {
						rule.count = 0
					}
				}
			}
		}
	}()
}
