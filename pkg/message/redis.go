package message

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

type RedisPubSub struct {
	client *redis.Client
}

func NewRedisPubSub() *RedisPubSub {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6385",
	})
	return &RedisPubSub{client: client}
}

func (r *RedisPubSub) Init(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisPubSub) Publish(ctx context.Context, to string, data interface{}) error {
	return r.client.Publish(ctx, to, data).Err()
}

func (r *RedisPubSub) Subscribe(ctx context.Context, from string, cb SubscribeHandler) (UnSubscriber, error) {
	sub := r.client.Subscribe(ctx, from)
	go func() {
		ctx := context.Background()
		for {
			m, err := sub.ReceiveMessage(ctx)
			if err != nil {
				cb(&Message{Error: err})
				if errors.Is(err, redis.ErrClosed) {
					break
				}
			}
			cb(&Message{Data: []byte(m.Payload)})
		}
	}()
	return &RedisUnSubscriber{sub: sub}, nil
}

func (r *RedisPubSub) Close() {
	r.client.Close()
}

type RedisUnSubscriber struct {
	sub *redis.PubSub
}

func (rs *RedisUnSubscriber) UnSubscribe() error {
	return rs.sub.Close()
}
