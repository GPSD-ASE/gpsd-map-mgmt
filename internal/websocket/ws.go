package websocket

import (
	"log"
	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()
	clients[conn] = true

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			break
		}
		broadcast <- string(msg)
		conn.WriteMessage(messageType, msg)
	}
}
