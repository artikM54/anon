package main

import (
	"anonymous_chat/internal/config"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

// for start "CONFIG_PATH="/home/cat/test_area/go/anonymous_chat" go run main.go"
func main() {
	config.MustLoad()
	setupRoutes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRoutes() {
	http.HandleFunc("/ws", wsHandler)
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	clientQueue []*Client
)

type Client struct {
	conn *websocket.Conn
}

func sendMessage(client *Client, message []byte) error {
	return client.conn.WriteMessage(websocket.TextMessage, message)
}

func chooseRandomPair() (*Client, *Client, error) {
	if len(clientQueue) < 2 {
		return nil, nil, errors.New("not enough clientQueue to pair")
	}

	idx1 := rand.Intn(len(clientQueue))
	idx2 := rand.Intn(len(clientQueue))

	for idx2 == idx1 {
		idx2 = rand.Intn(len(clientQueue))
	}

	return clientQueue[idx1], clientQueue[idx2], nil
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	defer conn.Close()

	client := &Client{conn: conn}
	clientQueue = append(clientQueue, client)

	bindClients()
}

func bindClients() {
	for {
		if len(clientQueue) >= 2 {
			fmt.Println("There are two users")

			client1, client2, err := chooseRandomPair()
			if err != nil {
				log.Println("Error choosing random pair:", err)
				continue
			}

			clientQueue = removeClientFromSlice(clientQueue, client1)
			clientQueue = removeClientFromSlice(clientQueue, client2)

			go handleStreamMessages(client1, client2)
			go handleStreamMessages(client2, client1)
		}
	}
}

func handleStreamMessages(client1 *Client, client2 *Client) {
	for {
		_, message, err := client1.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from client 1:", err)
			return
		}

		if err := sendMessage(client2, message); err != nil {
			log.Println("Error sending message to client 2:", err)
			return
		}
	}
}

func removeClientFromSlice(slice []*Client, client *Client) []*Client {
	for i, c := range slice {
		if c == client {
			return append(slice[:i], slice[i+1:]...)
		}
	}

	return slice
}
