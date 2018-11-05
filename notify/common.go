package notify

type Notify interface {
	Send(desc string, content ...string)
}