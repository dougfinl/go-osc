package osc

import (
	"fmt"
	"net"
)

/*
UDPServer provides functionality to receive OSC messages over UDP.
*/
type UDPServer struct {
	localAddr  *net.UDPAddr
	dispatcher *messageDispatcher
}

/*
NewUDPServer creates a UDP OSC server (for receiving OSC packets).
*/
func NewUDPServer(ip string, port int) (*UDPServer, error) {
	server := &UDPServer{dispatcher: &messageDispatcher{}}

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
Listen starts the server listening for OSC packets.
*/
func (s *UDPServer) Listen() error {
	listener, err := net.ListenUDP("udp", s.localAddr)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		// Read a datagram into the buffer
		buf := make([]byte, 65535)
		n, err := listener.Read(buf)
		if err != nil {
			return err
		}

		go s.handleIncomingData(buf[:n])
	}
}

/*
Handle adds an address pattern handler to the server. If the addressPattern is invalid an error is returned.
*/
func (s *UDPServer) Handle(addressPattern string, function func(*Message)) error {
	err := s.dispatcher.addHandler(addressPattern, function)

	return err
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
		s.dispatcher.dispatch(p.(*Message))
	case *Bundle:
		fmt.Println("ERROR server does not yet handle bundles")
	}
}
