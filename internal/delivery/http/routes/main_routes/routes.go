package main_routes

import (
	"net/http"
)

func SetupMainRoutes() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
}
