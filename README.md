# Server-Client Communication Project

This project implements a robust server and client system in Go, supporting TCP, UDP unicast, and UDP multicast communication. The system allows multiple clients to connect, authenticate, and send periodic time messages to the server. The server confirms receipt by sending ACK messages back to the client, with special handling for UDP multicast.

## Features

### Server
- Thread-safe singleton design.
- Supports TCP, UDP unicast, and UDP multicast protocols.
- Authenticates clients using username and password.
- Receives periodic time messages from clients.
- Sends acknowledgment messages (ACK) for received time messages (TCP and UDP unicast).
- For UDP multicast, ACKs are sent but may not always be received by clients due to network/multicast limitations.

### Client
- Connects to the server over TCP, UDP unicast, or UDP multicast.
- Authenticates with username and password (required for all protocols).
- Sends local time to the server every 7 seconds.
- Retries sending the time message if no ACK is received within 2 seconds, up to 5 attempts.
- For UDP multicast, missing ACKs are logged as warnings (not fatal), as multicast ACK delivery is unreliable.

## UDP Multicast Note
UDP multicast is inherently unreliable for two-way communication. While the server sends ACKs for multicast messages, clients may not always receive them due to:
- Network stack limitations
- Multicast group membership
- OS and router configurations
- The nature of multicast (one-to-many, not always routable back to sender)

**In this project, UDP multicast clients log a warning if an ACK is not received, but continue operation. This is expected and does not indicate a bug.**

## Test
A test scenario is included that spins up one server and four simultaneous clients:
- **Client A:** TCP
- **Client B:** UDP unicast
- **Client C:** UDP multicast
- **Client D:** UDP multicast

Each client authenticates and sends time messages to the server every 7 seconds. The server logs all received messages and authentication events.

## Getting Started

### Prerequisites
- Go 1.20+ installed
- Network connectivity between server and clients (local or LAN)
- For UDP multicast: run clients and server on the same subnet for best results

### Installation
Clone the repository:
```sh
git clone <repository_url>
cd server-client
```

### Running
To run the test scenario:
```sh
go run main.go
```

You should see logs for server startup, client connections, authentication, and periodic time/ACK messages. UDP multicast clients may log warnings about missing ACKs—this is normal.

## File Structure
- `main.go` — Entry point, starts server and clients for testing
- `server/` — Server implementation
- `client/` — Client implementation
- `protocol/` — Message definitions and serialization helpers
