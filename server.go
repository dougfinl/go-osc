package osc

import (
	"encoding/binary"
	"fmt"
	"net"
)

const udpReadBufSize = 4096

/*
Server provides functionality to receive OSC messages over UDP or TCP.
*/
type Server interface {
	SetLocalAddr(ip string, port int) error
	StartListening() error
	Handle(addressPattern string, fn MessageHandleFunc) error
}

/*
UDPServer provides functionality to receive OSC messages over UDP.
*/
type UDPServer struct {
	localAddr *net.UDPAddr

	AddressSpace
}

// Compile-time check to ensure UDPServer implements the Server interface.
var _ Server = &UDPServer{}

/*
NewUDPServer creates a UDP OSC server (for receiving OSC packets).
*/
func NewUDPServer(ip string, port int) (Server, error) {
	server := &UDPServer{}

	err := server.SetLocalAddr(ip, port)
	if err != nil {
		return nil, err
	}

	return server, nil
}

/*
SetLocalAddr sets the local address and port that the UDP server will listen upon.
*/
func (s *UDPServer) SetLocalAddr(ip string, port int) error {
	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	s.localAddr = localAddr

	return nil
}

/*
StartListening starts the server listening for OSC packets.
*/
func (s *UDPServer) StartListening() error {
	conn, err := net.ListenUDP("udp", s.localAddr)
	if err != nil {
		return err
	}

	// defer conn.Close()

	go s.listen(conn)

	return nil
}

func (s *UDPServer) listen(conn net.Conn) {
	for {
		// Read a datagram into the buffer
		buf := make([]byte, udpReadBufSize)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		go s.handleIncomingData(buf[:n])
	}
}

/*
handleIncomingData attempts to decode and dispatch the incoming OSC packet. If the data is not a valid OSC packet, it is silently ignored.
*/
func (s *UDPServer) handleIncomingData(data []byte) {
	p, err := decodePacket(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch p.(type) {
	case *Message:
		s.AddressSpace.Dispatch(p.(*Message))
	case *Bundle:
		fmt.Println("ERROR server does not yet handle bundles")
	}
}

/*
TCPServer provides functionality to receive OSC messages over TCP.
*/
type TCPServer struct {
	localAddr *net.TCPAddr

	AddressSpace
}

// Compile-time check to ensure TCPServer implements the Server interface.
var _ Server = &TCPServer{}

/*
NewTCPServer creates a TCP OSC server (for receiving OSC packets).
*/
func NewTCPServer(ip string, port int) (Server, error) {
	server := &TCPServer{}

	err := server.SetLocalAddr(ip, port)
	if err != nil {
		return nil, err
	}

	return server, nil
}

/*
SetLocalAddr sets the local address and port that the TCP server will listen upon.
*/
func (s *TCPServer) SetLocalAddr(ip string, port int) error {
	localAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	s.localAddr = localAddr

	return nil
}

/*
StartListening starts the server listening for incoming TCP connections.
*/
func (s *TCPServer) StartListening() error {
	listener, err := net.ListenTCP("tcp", s.localAddr)
	if err != nil {
		return err
	}

	defer listener.Close()

	go s.listen(listener)

	return nil
}

func (s *TCPServer) listen(listener net.Listener) {
	for {
		conn, err := listener.Accept()

		// Read a datagram into the buffer
		buf := make([]byte, udpReadBufSize)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		go s.handleIncomingData(buf[:n])
	}
}

/*
handleIncomingData attempts to decode and dispatch the incoming OSC packet. If the data is not a valid OSC packet encoded with a packet length header (OSC 1.0), it is silently ignored.
*/
func (s *TCPServer) handleIncomingData(data []byte) {
	// First four bytes should be the data length
	lenP := binary.BigEndian.Uint32(data[:4])
	fmt.Print(lenP)

	p, err := decodePacket(data[len(data)-3:])
	if err != nil {
		fmt.Println(err)
		return
	}

	switch p.(type) {
	case *Message:
		s.AddressSpace.Dispatch(p.(*Message))
	case *Bundle:
		fmt.Println("ERROR server does not yet handle bundles")
	}
}
