package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"server-client/protocol"
	"time"
)

type Client struct {
	Protocol string // "tcp", "udp-unicast", "udp-multicast"
	Addr     string
	Username string
	Password string
}

// Connect establishes a connection to the server (temp)
func (c *Client) Connect() error {
	if c.Protocol == "tcp" {
		return c.connectTCP()
	}
	return fmt.Errorf("protocol %s not implemented yet", c.Protocol)
}

func (c *Client) connectTCP() error {
	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return fmt.Errorf("tcp dial failed: %w", err)
	}
	fmt.Printf("[%s client] connected to %s\n", c.Protocol, c.Addr)

	auth := protocol.AuthMessage{Username: c.Username, Password: c.Password}
	authBytes, _ := json.Marshal(auth)
	authBytes = append(authBytes, '\n')
	_, _ = conn.Write(authBytes)

	tmsg := protocol.TimeMessage{Timestamp: time.Now()}
	tBytes, _ := json.Marshal(tmsg)
	tBytes = append(tBytes, '\n')
	_, _ = conn.Write(tBytes)

	reader := bufio.NewReader(conn)
	ackBytes, _ := reader.ReadBytes('\n')
	var ack protocol.AckMessage
	_ = json.Unmarshal(ackBytes, &ack)
	fmt.Println("[CLIENT] received ack:", ack.Status)

	return nil
}

// SendAuth sends authentication message (username + password) (temp)
func (c *Client) SendAuth() error {
	fmt.Printf("[%s client] Sending auth (user=%s)\n", c.Protocol, c.Username)
	return nil
}

// SendTime sends current time to server (temp)
func (c *Client) SendTime() error {
	fmt.Printf("[%s client] Sending time message\n", c.Protocol)
	return nil
}
