package osc

import (
	"encoding"
	"fmt"
)

/*
Packet represents and encodable OSC packet.
*/
type Packet interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	fmt.Stringer
}
