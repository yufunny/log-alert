package notify

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"strconv"
	"strings"
)

type MailNotify struct {
	Url       string
	Receivers []string
}

func getMailNotify (url string, receivers []string) (Notify, error)   {
	return &MailNotify{
		Url:url,
		Receivers:receivers,
	}, nil
}

func (x *MailNotify) Send(receivers []string,desc string, content ...string) {
	p := strings.Split(x.Url, "|")
	if len(p) != 4 {
		logrus.Errorf("mail config error")
		return
	}
	address := p[0]
	password := p[1]
	smtp := p[2]
	port, err:= strconv.Atoi(p[3])
	if err != nil {
		logrus.Errorf("mail port config error:%s", err.Error())
		return
	}

	m := gomail.NewMessage()
	// 发件人
	m.SetAddressHeader("From", address, "notice")
	// 收件人

	if len(receivers) == 0 {
		receivers = x.Receivers
	}

	m.SetHeader("To", receivers...)
	// 主题
	m.SetHeader("Subject", fmt.Sprintf("[%s]日志监控警报", desc))
	body := "日志内容:"
	for _, str := range content {
		body = body + "\n" + str
	}
	// 发送的body体
	m.SetBody("text/plain", body)
	d := gomail.NewDialer(smtp, port, address, password)
	if err := d.DialAndSend(m); err != nil {
		logrus.Errorf("[mail notify]error:%s", err.Error())
	} else {
		logrus.Info("mail send success...")
	}
}
