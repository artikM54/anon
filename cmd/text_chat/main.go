package main

import (
	"anonymous_chat/internal/config"
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	chatService "anonymous_chat/internal/services/chat"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// for start "CONFIG_PATH="/home/cat/test_area/go/anonymous_chat" go run main.go"
func main() {
	config.MustLoad()
	setupRoutes()

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRoutes() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
}

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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("try to connect")
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	// TODO: implement closing connects
	// defer conn.Close()

	user := userModel.NewUser(conn)
	chatService.AddUserToQueue(user)

	// TODO: move to a queue handler
	if chatService.GetCountUsersIntoQueue() >= 2 {
		chatService.BindClients()
	} else {
		message := messageModel.NewMessage(
			"system",
			string("Нет свободных участников, пожалуйста, дождитесь свободного участника"),
		)

		if err := chatService.SendMessage(user, message); err != nil {
			log.Println("Error sending message to client 2: ", err)
			return
		}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Тест"))
}
