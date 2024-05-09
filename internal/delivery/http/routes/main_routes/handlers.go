package main_routes

import (
	userService "anonymous_chat/internal/services/user"
	chatService "anonymous_chat/internal/services/chat"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins
			return true
		},
	}
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Тест"))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("try to connect at", time.Now().Format("2006-01-02 15:04:05"))
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	user := userService.NewUser(conn)
	// fmt.Println(*user)
	// fmt.Println(user.Hash)
	c := chatService.NewChatService()
	c.AddUserAndTryUpChat(user)
}
