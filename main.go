package main

import (
	"fmt"
	"server-client/client"
	"server-client/server"
	"time"
)

func main() {
	srv := server.GetServer()
	go func() {
		if err := srv.StartTCP(":9000"); err != nil {
			fmt.Println("TCP server error:", err)
		}
	}()
	go func() {
		if err := srv.StartUDPUnicast(":9001"); err != nil {
			fmt.Println("UDP unicast server error:", err)
		}
	}()
	go func() {
		if err := srv.StartUDPMulticast("224.0.0.1:9002"); err != nil {
			fmt.Println("UDP multicast server error:", err)
		}
	}()

	time.Sleep(500 * time.Millisecond)

	clients := []*client.Client{
		{Protocol: "tcp", Addr: "localhost:9000", Username: "alice", Password: "secret"},           //Client A
		{Protocol: "udp-unicast", Addr: "localhost:9001", Username: "bob", Password: "secret"},     //Client B
		{Protocol: "udp-multicast", Addr: "224.0.0.1:9002", Username: "chris", Password: "secret"}, //Client C
		{Protocol: "udp-multicast", Addr: "224.0.0.1:9002", Username: "dick", Password: "secret"},  //Client D
	}

	for _, c := range clients {
		c := c
		if err := c.Connect(); err != nil {
			fmt.Printf("[%s client] failed to connect: %v\n", c.Protocol, err)
			continue
		}

		go func(cli *client.Client) {
			for {
				if !cli.SendTime() {
					break
				}
				time.Sleep(7 * time.Second)
			}
		}(c)
	}

	select {}
}
