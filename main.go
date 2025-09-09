package main

import (
	"server-client/client"
	"server-client/server"
	"time"
)

func main() {
	// Start server
	srv := server.GetServer()
	go srv.StartTCP(":9000")
	go srv.StartUDPUnicast(":9001")
	go srv.StartUDPMulticast("224.0.0.1:9002")

	// Client A: TCP
	clientA := client.Client{
		Protocol: "tcp",
		Addr:     "localhost:9000",
		Username: "alice",
		Password: "secret",
	}
	if err := clientA.Connect(); err != nil {
		panic(err)
	}

	// Client B: UDP unicast
	clientB := client.Client{
		Protocol: "udp-unicast",
		Addr:     "localhost:9001",
		Username: "bob",
		Password: "secret",
	}
	if err := clientB.Connect(); err != nil {
		panic(err)
	}

	// Client C: UDP multicast
	clientC := client.Client{
		Protocol: "udp-multicast",
		Addr:     "224.0.0.1:9002",
		Username: "chris",
		Password: "secret",
	}
	if err := clientC.Connect(); err != nil {
		panic(err)
	}

	// Client D: UDP multicast
	clientD := client.Client{
		Protocol: "udp-multicast",
		Addr:     "224.0.0.1:9002",
		Username: "dick",
		Password: "secret",
	}
	if err := clientD.Connect(); err != nil {
		panic(err)
	}

	// Run them all in background loops
	go func() {
		for {
			clientA.SendTime()
			time.Sleep(7 * time.Second)
		}
	}()
	go func() {
		for {
			clientB.SendTime()
			time.Sleep(7 * time.Second)
		}
	}()
	go func() {
		for {
			clientC.SendTime()
			time.Sleep(7 * time.Second)
		}
	}()
	go func() {
		for {
			clientD.SendTime()
			time.Sleep(7 * time.Second)
		}
	}()

	select {}
}
