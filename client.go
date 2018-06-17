package osc

import (
	"fmt"
	"net"
)

/*
UDPClient provides functionality to send OSC messages over UDP.
*/
type UDPClient struct {
	addr      *net.UDPAddr
	localAddr *net.UDPAddr
}

/*
NewUDPClient creates a new UDP OSC client (for sending OSC packets).
*/
func NewUDPClient(ip string, port int) (*UDPClient, error) {
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
Send sends an OSC packet (message or bundle) from this client.
*/
func (c *UDPClient) Send(p Packet) error {
	conn, err := net.DialUDP("udp", c.localAddr, c.addr)
	if err != nil {
		return err
	}

	defer conn.Close()

	data, err := p.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}
