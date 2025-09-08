package main

import "server-client/server"

func main() {
	srv := server.GetServer()

	go srv.StartTCP(":9000")
	go srv.StartUDPUnicast(":9001")
	go srv.StartUDPMulticast(":9002")

	select {}
}
