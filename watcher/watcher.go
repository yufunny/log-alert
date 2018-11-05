package watcher

import (
	"github.com/papertrail/go-tail/follower"
	"io"
	"regexp"
	"fmt"
	"os"
	"github.com/yufunny/log-alert/notify"
	"time"
	"github.com/sirupsen/logrus"
)

type Watcher struct {
	File string
	Rule string
	Desc string
	Duration time.Duration
	Times int
	Interval time.Duration
	Notifier notify.Notify
	Count int
	Sent bool
	Text []string
}

func (w *Watcher) Watch() {
	t, _ := follower.New(w.File, follower.Config{
		Whence: io.SeekEnd,
		Offset: 0,
		Reopen: true,
	})

	go w.tick()
	reg := regexp.MustCompile(w.Rule)

	for line := range t.Lines() {
		if reg.Match(line.Bytes()) {
			logrus.Debugf("%s:%s", w.Desc, string(line.Bytes()))
			w.Text = append(w.Text, string(line.Bytes()))
			if len(w.Text) >= w.Times && !w.Sent {
				w.Notifier.Send(w.Desc, w.Text...)
				if w.Interval.Nanoseconds() > 0 {
					w.Sent = true
				}
				w.Text = make([]string, 0)
			}
		}
	}

	if t.Err() != nil {
		fmt.Fprintln(os.Stderr, t.Err())
	}
}

func (w *Watcher) tick ()  {
	if w.Duration.Nanoseconds() > 0  {
		countTicker := time.NewTicker(w.Duration)
		go func() {
			for {
				select {
				case <-countTicker.C:
					w.Count = 0
				}
			}

		}()
	}
	if w.Interval.Nanoseconds() > 0 {
		sentTicker := time.NewTicker(w.Interval)
		go func() {
			for {
				select {
				case <-sentTicker.C:
					w.Sent = false
				}
			}

		}()
	}
}
