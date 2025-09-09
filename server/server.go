package server

import (
	"bufio"
	"fmt"
	"net"
	"server-client/protocol"
	"sync"
)

type Server struct {
	mu sync.Mutex
}

var (
	instance *Server
	once     sync.Once
)

var udpClients = struct {
	sync.Mutex
	authed map[string]bool
}{authed: make(map[string]bool)}

func GetServer() *Server {
	once.Do(func() {
		instance = &Server{}
	})
	return instance
}

func (s *Server) StartTCP(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("[SERVER-TCP] TCP listen error: %w", err)
	}
	fmt.Println("[SERVER-TCP] TCP listening on", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[SERVER-TCP] accept error:", err)
			continue
		}
		go s.handleTCPConn(conn)
	}
}

func (s *Server) handleTCPConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	authBytes, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("[SERVER-TCP] failed to read auth:", err)
		return
	}

	var auth protocol.AuthMessage
	if err := protocol.Decode(authBytes, &auth); err != nil {
		fmt.Println("[SERVER-TCP] invalid auth json:", err)
		return
	}

	if !s.validateUser(auth.Username, auth.Password) {
		fmt.Println("[SERVER-TCP] invalid credentials for", auth.Username)
		return
	}
	fmt.Println("[SERVER-TCP] client authenticated:", auth.Username)

	ack := protocol.AckMessage{Status: "ACK_AUTH"}
	ackBytes, _ := protocol.Encode(ack)
	if _, err := conn.Write(ackBytes); err != nil {
		fmt.Println("[SERVER-TCP] failed to send ACK_AUTH:", err)
		return
	}

	for {
		timeBytes, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("[SERVER-TCP] connection closed:", err)
			return
		}

		var tmsg protocol.TimeMessage
		if err := protocol.Decode(timeBytes, &tmsg); err != nil {
			fmt.Println("[SERVER-TCP] invalid time message:", err)
			continue
		}

		fmt.Printf("[SERVER-TCP] received time from %s: %v\n", auth.Username, tmsg.Timestamp)

		ack := protocol.AckMessage{Status: "ACK_TIME"}
		ackBytes, _ := protocol.Encode(ack)
		if _, err := conn.Write(ackBytes); err != nil {
			fmt.Println("[SERVER-TCP] failed to send ACK_TIME:", err)
			return
		}
	}
}

func (s *Server) StartUDPUnicast(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("[SERVER-UDP] UDP unicast server on", addr)

	return s.handleUDPUnicastConn(conn)
}

func (s *Server) handleUDPUnicastConn(conn *net.UDPConn) error {
	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		var msg map[string]any
		if err := protocol.Decode(buf[:n], &msg); err != nil {
			fmt.Println("[SERVER-UDP] decode error:", err)
			continue
		}

		addrKey := clientAddr.String()

		if u, ok := msg["username"]; ok {
			if s.validateUser(fmt.Sprint(msg["username"]), fmt.Sprint(msg["password"])) {
				udpClients.Lock()
				udpClients.authed[addrKey] = true
				udpClients.Unlock()
				fmt.Printf("[SERVER-UDP] client authenticated: %s (addr=%s)\n", u, clientAddr)
				ack := protocol.AckMessage{Status: "ACK_AUTH"}
				data, _ := protocol.EncodeUDP(ack)
				conn.WriteToUDP(data, clientAddr)
			} else {
				fmt.Printf("[SERVER-UDP] invalid credentials for %s\n", u)
			}
			continue
		}

		if ts, ok := msg["timestamp"]; ok {
			udpClients.Lock()
			authed := udpClients.authed[addrKey]
			udpClients.Unlock()
			if !authed {
				fmt.Printf("[SERVER-UDP] ignoring time from unauthenticated client %s\n", addrKey)
				continue
			}

			fmt.Printf("[SERVER-UDP] received time from %s: %v\n", clientAddr, ts)
			ack := protocol.AckMessage{Status: "ACK_TIME"}
			data, _ := protocol.EncodeUDP(ack)
			conn.WriteToUDP(data, clientAddr)
		}
	}
}

func (s *Server) StartUDPMulticast(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenMulticastUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("[SERVER-MULTICAST] UDP multicast server on", addr)

	conn.SetReadBuffer(1024)
	return s.handleUDPMulticastConn(conn)
}

func (s *Server) handleUDPMulticastConn(conn *net.UDPConn) error {
	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		var msg map[string]any
		if err := protocol.Decode(buf[:n], &msg); err != nil {
			fmt.Println("[SERVER-MULTICAST] decode error:", err)
			continue
		}

		addrKey := clientAddr.String()

		if u, ok := msg["username"]; ok {
			if s.validateUser(fmt.Sprint(msg["username"]), fmt.Sprint(msg["password"])) {
				udpClients.Lock()
				udpClients.authed[addrKey] = true
				udpClients.Unlock()
				fmt.Printf("[SERVER-MULTICAST] client authenticated: %s (addr=%s)\n", u, clientAddr)
				ack := protocol.AckMessage{Status: "ACK_AUTH"}
				data, _ := protocol.EncodeUDP(ack)
				conn.WriteToUDP(data, clientAddr)
			} else {
				fmt.Printf("[SERVER-MULTICAST] invalid credentials for %s\n", u)
			}
			continue
		}

		if ts, ok := msg["timestamp"]; ok {
			udpClients.Lock()
			authed := udpClients.authed[addrKey]
			udpClients.Unlock()
			if !authed {
				fmt.Printf("[SERVER-MULTICAST] ignoring time from unauthenticated client %s\n", addrKey)
				continue
			}

			fmt.Printf("[SERVER-MULTICAST] received time from %s: %v\n", clientAddr, ts)
			ack := protocol.AckMessage{Status: "ACK_TIME"}
			data, _ := protocol.EncodeUDP(ack)
			conn.WriteToUDP(data, clientAddr)
		}
	}
}

func (s *Server) validateUser(username, password string) bool {
	return username != "" && password != ""
}
