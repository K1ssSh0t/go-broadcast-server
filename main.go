// Package main implements a WebSocket-based broadcast server and client
// that allows multiple clients to connect and exchange messages in real-time.
//
// Usage:
//
//	Start server: broadcast-server start [-port PORT]
//	Connect client: broadcast-server connect [-host HOST] [-port PORT]
//
// Example:
//
//	Start server on default port (8080):
//	  $ broadcast-server start
//
//	Start server on custom port:
//	  $ broadcast-server start -port 9000
//
//	Connect client to local server:
//	  $ broadcast-server connect
//
//	Connect client to remote server:
//	  $ broadcast-server connect -host example.com -port 9000
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// defaultPort is the default port number for both server and client connections
const defaultPort = 8080

// main is the entry point of the application that handles command-line arguments
// and starts either a server or client based on the provided subcommand
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: broadcast-server [start|connect]")
		os.Exit(1)
	}

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startPort := startCmd.Int("port", defaultPort, "Port to start the server on")

	connectCmd := flag.NewFlagSet("connect", flag.ExitOnError)
	connectPort := connectCmd.Int("port", defaultPort, "Port to connect to")
	connectHost := connectCmd.String("host", "localhost", "Host to connect to")

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		log.Printf("Starting server on port %d...\n", *startPort)
		startServer(*startPort)
	case "connect":
		connectCmd.Parse(os.Args[2:])
		log.Printf("Connecting to %s:%d...\n", *connectHost, *connectPort)
		startClient(*connectHost, *connectPort)
	default:
		fmt.Println("Expected 'start' or 'connect' subcommands")
		os.Exit(1)
	}
}