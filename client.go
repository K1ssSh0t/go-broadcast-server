// Package main provides the client implementation for the WebSocket broadcast system.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/fasthttp/websocket"
)

// startClient initializes and runs a WebSocket client that connects to the broadcast server
// It handles user input from the console and displays received messages from the server
// The client can be terminated using Ctrl+C
//
// Parameters:
//   - host: The hostname or IP address of the broadcast server
//   - port: The port number where the server is listening
//
// The client performs the following operations:
//   1. Establishes a WebSocket connection to the server
//   2. Starts a goroutine to handle incoming messages
//   3. Reads user input from the console and sends it to the server
//   4. Handles graceful shutdown on interrupt signal (Ctrl+C)
func startClient(host string, port int) {
	// Channel to handle interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Connect to the WebSocket server
	url := fmt.Sprintf("ws://%s:%d/ws", host, port)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Channel to signal when the read loop is done
	done := make(chan struct{})

	// Start goroutine to read messages from server
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("Received: %s", message)
		}
	}()

	// Read input from console and send to server
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Connected to server. Type your messages (press Ctrl+C to quit):")
	
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and waiting for server to close the connection
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		default:
			if scanner.Scan() {
				message := scanner.Text()
				err := c.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		}
	}
}