package main

import (
	"server-client/client"
	"server-client/server"
)

func main() {
	// Start server
	srv := server.GetServer()
	go srv.StartTCP(":9000")
	go srv.StartUDPUnicast(":9001")
	go srv.StartUDPMulticast(":9002")

	// Start sample client
	clientA := client.Client{
		Protocol: "tcp",
		Addr:     "localhost:9000",
		Username: "alice",
		Password: "secret",
	}

	clientA.Connect()
	clientA.SendAuth()
	clientA.SendTime()

	select {}
}
