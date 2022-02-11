package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/key-hh/chat-example/pkg/message"
)

type userContextKey int

const (
	pingPeriod                = 60 * time.Second
	userKey    userContextKey = 0
)

type ChannelHandler struct {
	upgrader websocket.Upgrader
	PubSub   message.PubSub
}

func (ch *ChannelHandler) Init() {
	ch.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
}

func (ch *ChannelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn, err := ch.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("conn err %v", err)
		return
	}
	defer conn.Close()

	roomID := r.URL.Query().Get("room")
	userID := r.FormValue("user")

	ctx = context.WithValue(ctx, userKey, userID)

	channel := Channel{
		wc:     conn,
		recvCh: make(chan []byte, 100),
		ticker: time.NewTicker(pingPeriod),
		user:   userID,
		room:   roomID,
	}

	unSubscriber, err := ch.PubSub.Subscribe(ctx, roomID, func(m *message.Message) {
		if m.Error != nil {
			log.Printf("Subscriber err %v", m.Error)
			return
		}
		channel.recvCh <- m.Data
	})
	if err != nil {
		log.Printf("conn err %v", err)
		return
	}
	defer unSubscriber.UnSubscribe()

	channel.pub = func(ctx context.Context, room string, msg string) error {
		return ch.PubSub.Publish(ctx, room, msg)
	}

	go channel.Sender(ctx)
	channel.Receiver(ctx)
}

type Channel struct {
	wc     *websocket.Conn
	recvCh chan []byte
	ticker *time.Ticker
	pub    func(context.Context, string, string) error
	user   string
	room   string
}

func (c *Channel) Sender(ctx context.Context) {
	user := ctx.Value(userKey).(string)

	for {
		_, msg, err := c.wc.ReadMessage()
		if err != nil {
			log.Printf("Sender read err %v", err)
			return
		}
		log.Println("Sender read message ", string(msg))

		if err := c.pub(ctx, c.room, makeMessage(user, msg)); err != nil {
			log.Printf("Sender err %v", err)
		}
	}
}

func makeMessage(user string, message []byte) string {
	return fmt.Sprintf("%s > %s", user, string(message))
}

func (c *Channel) Receiver(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("%s user %s room context done %v", ctx.Err())
			return
		case <-c.ticker.C:
			if err := c.wc.WriteMessage(websocket.TextMessage, []byte("ping message")); err != nil {
				log.Printf("ping write err %v", err)
			}
		case msg := <-c.recvCh:
			if err := c.wc.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("msg write err %v", err)
			}
		}
	}
}
