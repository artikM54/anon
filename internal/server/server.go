package server

import (
	"anonymous_chat/internal/delivery/http/routes"
	"fmt"
	"log"
	"net/http"
)

func MustLoad() {
	routes.SetupRoutes()

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
