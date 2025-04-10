// Package main provides the server implementation for the WebSocket broadcast system.
package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Client represents a connected WebSocket client
type Client struct {
	Conn *websocket.Conn
}

// BroadcastServer manages WebSocket connections and message broadcasting
// It maintains a list of connected clients and handles message distribution
type BroadcastServer struct {
	// clients holds all connected WebSocket clients
	clients map[*websocket.Conn]bool

	// register is a channel for registering new client connections
	register chan *websocket.Conn

	// unregister is a channel for removing disconnected clients
	unregister chan *websocket.Conn

	// broadcast is a channel for sending messages to all clients
	broadcast chan []byte

	// mutex provides thread-safe access to the clients map
	mutex sync.RWMutex
}

// NewBroadcastServer creates and initializes a new BroadcastServer instance
// It initializes all required channels and the clients map
func NewBroadcastServer() *BroadcastServer {
	return &BroadcastServer{
		clients:    make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		broadcast:  make(chan []byte),
	}
}

// handleWebSocket manages an individual WebSocket connection
// It handles client registration, message reading, and cleanup on disconnect
func (s *BroadcastServer) handleWebSocket(c *websocket.Conn) {
	// Register new client
	s.register <- c
	defer func() {
		s.unregister <- c
		c.Close()
	}()

	// Handle incoming messages
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			return
		}

		if messageType == websocket.TextMessage {
			s.broadcast <- message
		}
	}
}

// run starts the main broadcast server loop
// It processes registration, unregistration, and message broadcasting
// This method should be run in a separate goroutine
func (s *BroadcastServer) run() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			s.mutex.Unlock()
			log.Printf("Client connected. Total clients: %d", len(s.clients))

		case client := <-s.unregister:
			s.mutex.Lock()
			delete(s.clients, client)
			s.mutex.Unlock()
			log.Printf("Client disconnected. Total clients: %d", len(s.clients))

		case message := <-s.broadcast:
			s.mutex.RLock()
			for client := range s.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("error broadcasting message: %v", err)
					client.Close()
					delete(s.clients, client)
				}
			}
			s.mutex.RUnlock()
		}
	}
}

// startServer initializes and starts the Fiber web server with WebSocket support
// It creates routes for WebSocket connections and starts listening on the specified port
func startServer(port int) {
	app := fiber.New()
	server := NewBroadcastServer()

	// Start the broadcast handler
	go server.run()

	// WebSocket route
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		server.handleWebSocket(c)
	}))

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}