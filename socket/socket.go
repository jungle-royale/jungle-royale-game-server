package socket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var ConnectionChannel chan *websocket.Conn

func InitSocket() {

	ConnectionChannel = make(chan *websocket.Conn, 100)
	http.HandleFunc("/ws", handleWebSocket)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// http → websocket upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	ConnectionChannel <- conn
}
