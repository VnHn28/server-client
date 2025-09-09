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
	authBytes, _ := protocol.EncodeUDP(auth)
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
	authBytes, _ := protocol.EncodeUDP(auth)
	if _, err := conn.Write(authBytes); err != nil {
		return fmt.Errorf("failed to send auth: %w", err)
	}

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Printf("[%s client] warning: did not receive ACK (multicast ACKs may be unreliable): %v\n", c.Protocol, err)
		return nil
	}
	var ack protocol.AckMessage
	if err := protocol.Decode(buf[:n], &ack); err != nil {
		fmt.Printf("[%s client] warning: failed to decode ACK (multicast): %v\n", c.Protocol, err)
		return nil
	}
	fmt.Printf("[%s client] received ack: %s\n", c.Protocol, ack.Status)

	return nil
}

// SendTime sends current time to server and waits for ACK
// Retries up to 5 times if no ACK received within 2 seconds
func (c *Client) SendTime() {
	if c.conn == nil {
		fmt.Printf("[%s client] no connection, cannot send time\n", c.Protocol)
		return
	}

	tmsg := protocol.TimeMessage{Timestamp: time.Now()}
	var data []byte
	if c.Protocol == "udp-unicast" || c.Protocol == "udp-multicast" {
		data, _ = protocol.EncodeUDP(tmsg)
	} else {
		data, _ = protocol.Encode(tmsg)
	}

	maxRetries := 5
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if _, err := c.conn.Write(data); err != nil {
			fmt.Printf("[%s client] failed to send time: %v\n", c.Protocol, err)
			return
		}

		if err := c.conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
			fmt.Printf("[%s client] failed to set deadline: %v\n", c.Protocol, err)
			return
		}

		buf := make([]byte, 1024)
		n, _, err := c.readFromConn(buf)
		if err == nil {
			var ack protocol.AckMessage
			if decodeErr := protocol.Decode(buf[:n], &ack); decodeErr == nil {
				fmt.Printf("[%s client] received ack: %s\n", c.Protocol, ack.Status)
				return
			}
		}

		fmt.Printf("[%s client] attempt %d/%d: no ack received, retrying...\n", c.Protocol, attempt, maxRetries)
	}

	fmt.Printf("[%s client] no ack after %d attempts, closing connection\n", c.Protocol, maxRetries)
	c.conn.Close()
	c.conn = nil
}

func (c *Client) readFromConn(buf []byte) (int, net.Addr, error) {
	switch conn := c.conn.(type) {
	case *net.UDPConn:
		return conn.ReadFrom(buf)
	default: // TCP
		n, err := conn.Read(buf)
		return n, nil, err
	}
}
