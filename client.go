package osc

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

/*
Client represents an OSC client (UDP or TCP) that can send OSC packets to a remote host.
*/
type Client interface {
	SetAddr(ip string, port int) error
	SetLocalAddr(ip string, port int) error
	Connect() error
	Disconnect() error
	IsConnected() bool
	Send(p Packet) error
}

/*
UDPClient provides functionality to send OSC messages over UDP.
*/
type UDPClient struct {
	addr      *net.UDPAddr
	localAddr *net.UDPAddr
	conn      *net.UDPConn
	connected bool
}

// Compile-time check to ensure UDPClient implements the Client interface.
var _ Client = &UDPClient{}

/*
NewUDPClient creates a new UDP OSC client (for sending OSC packets).
*/
func NewUDPClient(ip string, port int) (Client, error) {
	client := &UDPClient{}

	err := client.SetAddr(ip, port)
	if err != nil {
		return nil, err
	}

	return client, nil
}

/*
SetAddr sets the destination address for packets send by this client.
*/
func (c *UDPClient) SetAddr(ip string, port int) error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	c.addr = addr

	return nil
}

/*
SetLocalAddr sets the local address for packets to be sent from by this client.
*/
func (c *UDPClient) SetLocalAddr(ip string, port int) error {
	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	c.localAddr = localAddr

	return nil
}

/*
Connect connects the client to the remote host.
*/
func (c *UDPClient) Connect() error {
	conn, err := net.DialUDP("udp", c.localAddr, c.addr)
	if err != nil {
		return err
	}

	c.conn = conn

	c.connected = true

	return nil
}

/*
Disconnect disconnects the client from the remote host.
*/
func (c *UDPClient) Disconnect() error {
	if c.IsConnected() {
		return c.conn.Close()
	}

	return nil
}

/*
IsConnected returns true if the client is connected to the remote host.
*/
func (c UDPClient) IsConnected() bool {
	return c.conn != nil && c.connected
}

/*
Send sends an OSC packet (message or bundle) from this client.
*/
func (c *UDPClient) Send(p Packet) error {
	if !c.IsConnected() {
		return fmt.Errorf("Client is not connected")
	}

	data, err := p.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = c.conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

/*
TCPClient provides functionality to stream OSC messages to a remote host.
It also contains an AddressSpace to handle responses over the TCP stream.
*/
type TCPClient struct {
	addr      *net.TCPAddr
	localAddr *net.TCPAddr
	conn      *net.TCPConn
	connected bool

	AddressSpace
}

// Compile-time check to ensure TCPClient implements the Client interface.
var _ Client = &TCPClient{}

/*
NewTCPClient creates a new TCP OSC client (for sending OSC packets).
*/
func NewTCPClient(ip string, port int) (Client, error) {
	client := &TCPClient{}

	err := client.SetAddr(ip, port)
	if err != nil {
		return nil, err
	}

	return client, nil
}

/*
SetAddr sets the destination address for this connection.
*/
func (c *TCPClient) SetAddr(ip string, port int) error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	c.addr = addr

	return nil
}

/*
SetLocalAddr sets the local address for packets to be sent from by this client.
*/
func (c *TCPClient) SetLocalAddr(ip string, port int) error {
	localAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	c.localAddr = localAddr

	return nil
}

/*
Connect connects the TCPClient to the remote host.
*/
func (c *TCPClient) Connect() error {
	conn, err := net.DialTCP("tcp", c.localAddr, c.addr)
	if err != nil {
		return err
	}

	go c.responseReaderLoop()

	c.conn = conn

	return nil
}

func (c *TCPClient) responseReaderLoop() {
	buf := make([]byte, 65535)
	reader := bufio.NewReader(c.conn)

	for {
		var count uint32
		err := binary.Read(reader, binary.BigEndian, &count)
		if err != nil {
			fmt.Println("WARNING found malformed packet")
			break
		}

		_, err = io.ReadFull(reader, buf[:int(count)])
		if err != nil {
			fmt.Println("WARNING found malformed packet")
			break
		}

		p, err := decodePacket(buf[:int(count)])
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch p.(type) {
		case *Message:
			c.AddressSpace.Dispatch(p.(*Message))
		case *Bundle:
			fmt.Println("ERROR bundles not yet supported")
		}
	}
}

/*
Disconnect closes the TCPClient's connection.
*/
func (c *TCPClient) Disconnect() error {
	return c.conn.Close()
}

/*
IsConnected returns true if the client is connected to the remote host.
*/
func (c TCPClient) IsConnected() bool {
	return c.conn != nil && c.connected
}

/*
Send sends an OSC packet (message or bundle) from this client.
*/
func (c *TCPClient) Send(p Packet) error {
	packetEnc, err := p.MarshalBinary()
	if err != nil {
		return err
	}

	// Count the data to be sent, and encode as uint32 (OSC 1.0)
	count := len(packetEnc)
	countEnc := make([]byte, 4)
	binary.BigEndian.PutUint32(countEnc, uint32(count))

	data := append(countEnc, packetEnc...)

	_, err = c.conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}
