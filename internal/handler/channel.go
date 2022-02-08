package handler

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type ChannelHandler struct {
	upgrader websocket.Upgrader
}

func (ch *ChannelHandler) Init() {
	ch.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
}

func (ch *ChannelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := ch.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("conn err %v", err)
		return
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read err %v", err)
			return
		}

		log.Println("server got message ", msgType, string(msg))

		msgSend := []byte("echo response:" + string(msg))
		if err := conn.WriteMessage(msgType, msgSend); err != nil {
			log.Printf("write err %v", err)
			return
		}
	}
}
