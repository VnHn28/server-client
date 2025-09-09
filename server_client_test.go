package main

import (
	"fmt"
	"server-client/client"
	"server-client/server"
	"sync"
	"testing"
	"time"
)

func TestServerClientCommunication(t *testing.T) {
	srv := server.GetServer()
	var wg sync.WaitGroup

	// Start server concurrently
	go func() {
		if err := srv.StartTCP(":9000"); err != nil {
			t.Errorf("TCP server error: %v", err)
		}
	}()
	go func() {
		if err := srv.StartUDPUnicast(":9001"); err != nil {
			t.Errorf("UDP unicast server error: %v", err)
		}
	}()
	go func() {
		if err := srv.StartUDPMulticast("224.0.0.1:9002"); err != nil {
			t.Errorf("UDP multicast server error: %v", err)
		}
	}()

	time.Sleep(500 * time.Millisecond)

	clients := []*client.Client{
		{Protocol: "tcp", Addr: "localhost:9000", Username: "alice", Password: "secret"},
		{Protocol: "udp-unicast", Addr: "localhost:9001", Username: "bob", Password: "secret"},
		{Protocol: "udp-multicast", Addr: "224.0.0.1:9002", Username: "chris", Password: "secret"},
		{Protocol: "udp-multicast", Addr: "224.0.0.1:9002", Username: "dick", Password: "secret"},
	}

	for _, c := range clients {
		wg.Add(1)
		c := c
		if err := c.Connect(); err != nil {
			t.Errorf("[%s client] failed to connect: %v", c.Protocol, err)
			wg.Done()
			continue
		}

		go func(cli *client.Client) {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				cli.SendTime()
				time.Sleep(7 * time.Second)
			}
			fmt.Printf("[%s client] completed test iterations\n", cli.Protocol)
		}(c)
	}

	wg.Wait()
}
