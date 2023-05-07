package sender

import (
	"context"
)

type Sender struct {
	notifySender  NotifySender
	eventReceiver EventReceiver
	logger        Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warning(msg string)
	Debug(msg string)
}

type NotifySender interface {
	Send(ctx context.Context, data interface{}) error
}

type EventReceiver interface {
	Read(ctx context.Context) ([]byte, error)
}

func NewSender(_receiver EventReceiver, _sender NotifySender, _logger Logger) *Sender {
	s := &Sender{
		eventReceiver: _receiver,
		notifySender:  _sender,
		logger:        _logger,
	}
	return s
}

func (s *Sender) Start(ctx context.Context) {
	go func() {
		for {
			msg, err := s.eventReceiver.Read(ctx)
			if err != nil {
				return
			}
			if err = s.notifySender.Send(ctx, string(msg)); err != nil {
				s.logger.Error("Ошибка отправки оповещения: " + err.Error())
			}
		}
	}()
}
