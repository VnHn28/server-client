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
	go srv.StartUDPMulticast(":9002")

	// Start TCP client
	clientA := client.Client{
		Protocol: "tcp",
		Addr:     "localhost:9000",
		Username: "alice",
		Password: "secret",
	}

	if err := clientA.Connect(); err != nil {
		panic(err)
	}
	clientA.SendAuth()

	for {
		clientA.SendTime()
		time.Sleep(7 * time.Second)
	}
}
