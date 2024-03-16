package main

import (
	"fmt"
	"net"
)

const listenPort = ":8080" // Replace with a different port if needed

func main() {
	listener, err := net.Listen("tcp", listenPort)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close() // Close listener on exit

	fmt.Println("Listening on port", listenPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn) // Handle connection concurrently
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close() // Close connection on exit

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading message:", err)
		return
	}

	message := string(buf[:n])
	fmt.Println("Received message:", message) // Log received message
}
