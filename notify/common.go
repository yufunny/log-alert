package notify

import (
	"sync"
	"fmt"
)

type Notify interface {
	Send(receivers []string, desc string, content ...string)
}

type notifierDriver func(url string, receivers []string) (Notify, error)

var drivers sync.Map

func init()  {
	register("mail", getMailNotify)
}

func register(scheme string, driver notifierDriver) error {
	_, loaded := drivers.LoadOrStore(scheme, driver)
	if loaded {
		return fmt.Errorf("register: notifier driver scheme '%s' alrealy exists", scheme)
	}

	return nil
}


func Open(schema, url string, receivers []string) (Notify, error) {
	driver, ok := drivers.Load(schema)
	if !ok {
		return nil, fmt.Errorf("open: notifier driver '%s' not exists", schema)
	}
	return driver.(notifierDriver)(url, receivers)
}
