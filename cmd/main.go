package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/ws", handleConnection)

	env := os.Getenv("ENVIRONMENT")
	webSocketPort := os.Getenv("WEB_SOCKET_PORT")
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", webSocketPort),
	}

	log.Println("Starting WebSocket server...")

	switch env {
	case "production":
		certFile := os.Getenv("CERT_FILE_PATH")
		keyFile := os.Getenv("KEY_FILE_PATH")

		server.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
			log.Fatal("TLS server error:", err)
		}
	case "development":
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Server error:", err)
		}
	default:
		log.Fatalf("Unknown environment: %s", env)
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket Connection")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Message reading error:", err)
			break
		}
		log.Printf("Incoming message: %s", message)

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Message writing error:", err)
			break
		}
		log.Printf("Message echoed back: %s", message)
	}
}
