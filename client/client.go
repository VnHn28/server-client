package client

import "fmt"

// Client represents a network client that can connect to the server using either TCP or UDP (unicast/multicast).
type Client struct {
	Protocol string // "tcp", "udp-unicast", "udp-multicast"
	Addr     string
	Username string
	Password string
}

// Connect establishes a connection to the server (temp)
func (c *Client) Connect() error {
	fmt.Printf("[%s client] Connecting to %s\n", c.Protocol, c.Addr)
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
