# WebSocket Broadcast Server

A real-time broadcast server implementation in Go using WebSocket protocol. This application allows multiple clients to connect to a central server and exchange messages in real-time, where messages sent by one client are broadcasted to all connected clients.

## Features

- WebSocket-based real-time communication
- Support for multiple concurrent client connections
- Command-line interface for both server and client modes
- Graceful handling of client connections and disconnections
- Thread-safe message broadcasting
- Clean shutdown handling

## Requirements

- Go 1.16 or higher
- [Fiber v2](https://github.com/gofiber/fiber)
- [Fiber WebSocket](https://github.com/gofiber/websocket)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-broadcast-server.git
   cd go-broadcast-server
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the application:
   ```bash
   go build -o broadcast-server
   ```

## Usage

### Starting the Server

To start the broadcast server on the default port (8080):
```bash
./broadcast-server start
```

To start the server on a specific port:
```bash
./broadcast-server start -port 9000
```

### Connecting Clients

To connect a client to a local server:
```bash
./broadcast-server connect
```

To connect to a remote server:
```bash
./broadcast-server connect -host example.com -port 9000
```

### Client Usage

Once connected:
1. Type your message and press Enter to send
2. Messages from other clients will appear automatically
3. Press Ctrl+C to disconnect

## Implementation Details

### Server Components

- **BroadcastServer**: Main server structure that manages client connections and message broadcasting
- **WebSocket Handler**: Manages individual client connections and message routing
- **Thread-safe Client Management**: Uses mutex to safely handle concurrent client operations
- **Channel-based Communication**: Utilizes Go channels for client registration and message broadcasting

### Client Components

- **Interactive Console**: Reads user input and displays received messages
- **Concurrent Message Handling**: Uses goroutines to handle sending and receiving messages
- **Graceful Shutdown**: Implements proper WebSocket connection closure

## Error Handling

The application includes comprehensive error handling for:
- Connection failures
- Message transmission errors
- Unexpected client disconnections
- Server shutdown
- Invalid command-line arguments

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.