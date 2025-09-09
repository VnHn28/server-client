package client

import (
	"bufio"
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

func (c *Client) Connect() error {
	switch c.Protocol {
	case "tcp":
		return c.connectTCP()
	case "udp-unicast":
		return c.connectUDPUnicast()
	case "udp-multicast":
		return c.connectUDPMulticast()
	default:
		return fmt.Errorf("unsupported protocol: %s", c.Protocol)
	}
}

func (c *Client) connectTCP() error {
	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return fmt.Errorf("tcp dial failed: %w", err)
	}
	c.conn = conn
	fmt.Printf("[%s client] connected to %s\n", c.Protocol, c.Addr)

	auth := protocol.AuthMessage{Username: c.Username, Password: c.Password}
	authBytes, _ := protocol.Encode(auth)
	if _, err := conn.Write(authBytes); err != nil {
		return fmt.Errorf("failed to send auth: %w", err)
	}

	reader := bufio.NewReader(conn)
	ackBytes, err := reader.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("failed to read ack: %w", err)
	}
	var ack protocol.AckMessage
	if err := protocol.Decode(ackBytes, &ack); err != nil {
		return fmt.Errorf("decode ack failed: %w", err)
	}
	fmt.Printf("[%s client] received ack: %s\n", c.Protocol, ack.Status)

	return nil
}

func (c *Client) connectUDPUnicast() error {
	conn, err := net.Dial("udp", c.Addr)
	if err != nil {
		return fmt.Errorf("udp-unicast dial failed: %w", err)
	}
	c.conn = conn
	fmt.Printf("[%s client] connected to %s\n", c.Protocol, c.Addr)

	auth := protocol.AuthMessage{Username: c.Username, Password: c.Password}
	authBytes, _ := protocol.Encode(auth)
	if _, err := conn.Write(authBytes); err != nil {
		return fmt.Errorf("failed to send auth: %w", err)
	}

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read ack: %w", err)
	}
	var ack protocol.AckMessage
	if err := protocol.Decode(buf[:n], &ack); err != nil {
		return fmt.Errorf("decode ack failed: %w", err)
	}
	fmt.Printf("[%s client] received ack: %s\n", c.Protocol, ack.Status)

	return nil
}

func (c *Client) connectUDPMulticast() error {
	laddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	if err != nil {
		return err
	}
	raddr, err := net.ResolveUDPAddr("udp", c.Addr)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		return fmt.Errorf("udp-multicast dial failed: %w", err)
	}
	c.conn = conn
	fmt.Printf("[%s client] connected to %s\n", c.Protocol, c.Addr)

	auth := protocol.AuthMessage{Username: c.Username, Password: c.Password}
	authBytes, _ := protocol.Encode(auth)
	if _, err := conn.Write(authBytes); err != nil {
		return fmt.Errorf("failed to send auth: %w", err)
	}

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return fmt.Errorf("failed to read ack: %w", err)
	}
	var ack protocol.AckMessage
	if err := protocol.Decode(buf[:n], &ack); err != nil {
		return fmt.Errorf("decode ack failed: %w", err)
	}
	fmt.Printf("[%s client] received ack: %s\n", c.Protocol, ack.Status)

	return nil
}

// SendTime sends current time to server and waits for ACK
// Retries up to 5 times if no ACK received within 2 seconds
func (c *Client) SendTime() error {
	if c.conn == nil {
		return fmt.Errorf("connection not established")
	}

	tmsg := protocol.TimeMessage{
		Timestamp: time.Now(),
	}
	data, err := protocol.Encode(tmsg)
	if err != nil {
		return err
	}

	for attempt := 1; attempt <= 5; attempt++ {
		if _, err := c.conn.Write(data); err != nil {
			return err
		}
		fmt.Printf("[%s client] Sending time message (attempt %d)\n", c.Protocol, attempt)

		c.conn.SetReadDeadline(time.Now().Add(2 * time.Second))

		var ack protocol.AckMessage
		buf := make([]byte, 1024)

		switch conn := c.conn.(type) {
		case *net.TCPConn:
			reader := bufio.NewReader(conn)
			raw, err := reader.ReadBytes('\n')
			if err == nil {
				if decodeErr := protocol.Decode(raw, &ack); decodeErr == nil {
					fmt.Printf("[%s client] received ack: %s\n", c.Protocol, ack.Status)
					return nil
				}
			}

		case *net.UDPConn:
			n, _, err := conn.ReadFromUDP(buf)
			if err == nil {
				if decodeErr := protocol.Decode(buf[:n], &ack); decodeErr == nil {
					fmt.Printf("[%s client] received ack: %s\n", c.Protocol, ack.Status)
					return nil
				}
			}
		}

		fmt.Printf("[%s client] no ack, retrying...\n", c.Protocol)
	}

	fmt.Printf("[%s client] no ack after 5 attempts, closing connection\n", c.Protocol)
	c.conn.Close()
	return fmt.Errorf("failed to receive ack after 5 attempts")
}
