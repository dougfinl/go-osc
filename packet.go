package osc

import "encoding"

/*
Packet represents and encodable OSC packet.
*/
type Packet interface {
	encoding.BinaryMarshaler
}
