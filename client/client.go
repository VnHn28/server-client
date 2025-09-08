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
	conn     net.Conn
}

// Connect establishes a connection to the server
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

// SendAuth sends authentication message
func (c *Client) SendAuth() error {
	auth := protocol.AuthMessage{
		Username: c.Username,
		Password: c.Password,
	}
	data, err := protocol.Encode(auth)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(data)
	if err != nil {
		return err
	}
	fmt.Printf("[%s client] Sending auth (user=%s)\n", c.Protocol, c.Username)
	return nil
}

// SendTime sends current time to server and waits for ACK
// Retries up to 5 times if no ACK received within 2 seconds
func (c *Client) SendTime() error {
	tmsg := protocol.TimeMessage{
		Timestamp: time.Now(),
	}
	data, err := protocol.Encode(tmsg)
	if err != nil {
		return err
	}

	for attempt := 1; attempt <= 5; attempt++ {
		// Send message
		_, err := c.conn.Write(data)
		if err != nil {
			return err
		}
		fmt.Printf("[%s client] Sending time message (attempt %d)\n", c.Protocol, attempt)

		// Wait for ACK with timeout
		c.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		reader := bufio.NewReader(c.conn)
		raw, err := reader.ReadBytes('\n')
		if err == nil {
			var ack protocol.AckMessage
			if decodeErr := protocol.Decode(raw, &ack); decodeErr == nil {
				fmt.Printf("[CLIENT] received ack: %s\n", ack.Status)
				return nil // Success
			}
		}

		fmt.Printf("[CLIENT] no ack, retrying...\n")
	}

	fmt.Printf("[CLIENT] no ack after 5 attempts, closing connection\n")
	c.conn.Close()
	return fmt.Errorf("failed to receive ack after 5 attempts")
}
