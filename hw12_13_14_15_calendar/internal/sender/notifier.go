package sender

import (
	"io"
	"strings"
)

type Notifier struct {
	out io.Writer
}

func NewNotifier(_out io.Writer) *Notifier {
	n := &Notifier{
		out: _out,
	}
	return n
}

func (n *Notifier) Send(msg string) error {
	_, err := io.Copy(n.out, strings.NewReader(msg))
	return err
}
