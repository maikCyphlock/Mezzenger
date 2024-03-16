package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins for now (adjust for production)
}

// Client represents a connected WebSocket client
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

var clients = make(map[*websocket.Conn]*Client) // Map to store connected clients
var mutex = &sync.Mutex{}                       // Mutex for thread-safe access to clients map

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading:", err)
		return
	}

	client := &Client{conn: conn, send: make(chan []byte)}
	defer func() {
		conn.Close()
	}()
	clients[conn] = client

	go handleClient(client)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			deleteClient(conn)
			break
		}
		log.Printf("Received message: %s (type: %d)", message, messageType)
		broadcastMessage(message) // Broadcast received message to all clients
	}
}

func handleClient(client *Client) {
	defer connClose(client.conn)
	for message := range client.send {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error sending message to client:", err)
			deleteClient(client.conn)
			return
		}
	}
}

func broadcastMessage(message []byte) {
	mutex.Lock()
	defer mutex.Unlock()

	for conn, client := range clients {
		select {
		case client.send <- message:
		default:
			log.Println("Client connection closed:", conn.RemoteAddr())
			deleteClient(conn)
		}
	}
}

func deleteClient(conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	if client, ok := clients[conn]; ok {
		delete(clients, conn)
		close(client.send)
	}
}

func connClose(conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	deleteClient(conn)
}

func main() {
	http.HandleFunc("/ws", handleWS)
	log.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
