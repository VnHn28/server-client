package server

import (
	"fmt"
	"sync"
)

type Server struct{}

var (
	instance *Server
	once     sync.Once
)

// GetServer returns the singleton instance of Server
func GetServer() *Server {
	once.Do(func() {
		instance = &Server{}
	})
	return instance
}

// StartTCP starts the TCP listener (temp)
func (s *Server) StartTCP(addr string) error {
	fmt.Println("Starting TCP server on", addr)
	return nil
}

// StartUDPUnicast starts the UDP unicast listener (temp)
func (s *Server) StartUDPUnicast(addr string) error {
	fmt.Println("Starting UDP unicast server on", addr)
	return nil
}

// StartUDPMulticast starts the UDP multicast listener (temp)
func (s *Server) StartUDPMulticast(addr string) error {
	fmt.Println("Starting UDP multicast server on", addr)
	return nil
}
