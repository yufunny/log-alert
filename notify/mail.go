package notify

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type MailNotify struct {
	Url       string
	Receivers []string
}

func (x *MailNotify) Send(desc string, content ...string) {
	var address, password, smtp string
	var port int
	fmt.Sscanf(x.Url, "%s | %s | %s | %d", &address, &password, &smtp, &port)

	m := gomail.NewMessage()
	// 发件人
	m.SetAddressHeader("From", address, "notice")
	// 收件人
	m.SetHeader("To", x.Receivers...)
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
	}
	fmt.Println("mail send success...")
}
