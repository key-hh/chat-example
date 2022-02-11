package main

import (
	"context"
	"github.com/key-hh/chat-example/pkg/message"
	"log"
	"net/http"

	"github.com/key-hh/chat-example/internal/handler"
)

const (
	addr = "0.0.0.0:8090"
)

func main() {
	ctx := context.Background()

	pubsub := message.NewRedisPubSub()
	if err := pubsub.Init(ctx); err != nil {
		log.Fatal(err)
	}
	defer pubsub.Close()

	ch := &handler.ChannelHandler{PubSub: pubsub}
	ch.Init()

	http.Handle("/channel", ch)

	log.Println("server is ready for ", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		if err == http.ErrServerClosed {
			log.Print("server closed ..")
		} else {
			log.Fatal(err)
		}
	}
}
