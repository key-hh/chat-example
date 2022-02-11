package message

import "context"

type Message struct {
	Data  []byte
	Error error
}

type SubscribeHandler = func(m *Message)

type PubSub interface {
	Publish(ctx context.Context, to string, data interface{}) error
	Subscribe(ctx context.Context, from string, cb SubscribeHandler) (UnSubscriber, error)
}

type UnSubscriber interface {
	UnSubscribe() error
}
