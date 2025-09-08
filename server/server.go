package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"server-client/protocol"
	"sync"
)

type Server struct {
	mu sync.Mutex
}

var (
	instance *Server
	once     sync.Once
)

// GetServer returns the singleton instance of Server
func GetServer() *Server {
	once.Do(func() {
		instance = &Server{}
	})
	return instance
}

// StartTCP starts the TCP listener and handles incoming connections
func (s *Server) StartTCP(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("TCP listen error: %w", err)
	}
	fmt.Println("[SERVER] TCP listening on", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[SERVER] accept error:", err)
			continue
		}
		go s.handleTCPConn(conn)
	}
}

func (s *Server) handleTCPConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	authBytes, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("[SERVER] failed to read auth:", err)
		return
	}

	var auth protocol.AuthMessage
	if err := json.Unmarshal(authBytes, &auth); err != nil {
		fmt.Println("[SERVER] invalid auth json:", err)
		return
	}

	if !s.validateUser(auth.Username, auth.Password) {
		fmt.Println("[SERVER] invalid credentials for", auth.Username)
		return
	}
	fmt.Println("[SERVER] client authenticated:", auth.Username)

	for {
		timeBytes, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("[SERVER] connection closed:", err)
			return
		}

		var tmsg protocol.TimeMessage
		if err := json.Unmarshal(timeBytes, &tmsg); err != nil {
			fmt.Println("[SERVER] invalid time message:", err)
			continue
		}

		fmt.Printf("[SERVER] received time from %s: %v\n", auth.Username, tmsg.Timestamp)

		// send ACK
		ack := protocol.AckMessage{Status: "OK"}
		ackBytes, _ := json.Marshal(ack)
		ackBytes = append(ackBytes, '\n')
		if _, err := conn.Write(ackBytes); err != nil {
			fmt.Println("[SERVER] failed to send ACK:", err)
			return
		}
	}
}

// StartUDPUnicast starts the UDP unicast listener (temp)
func (s *Server) StartUDPUnicast(addr string) error {
	fmt.Println("Starting UDP unicast server on", addr)
	return nil
}

// StartUDPMulticast starts the UDP multicast listener (temp)
func (s *Server) StartUDPMulticast(addr string) error {
	fmt.Println("Starting UDP multicast server on", addr)
	return nil
}

func (s *Server) validateUser(username, password string) bool {
	return username != "" && password != ""
}
