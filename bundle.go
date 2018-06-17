package osc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

var (
	// Encoded bundle identifier string
	bundleString = []byte{'#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00'}
)

/*
Bundle represents an OSC bundle, which contains a time tag and multiple child elements.
*/
type Bundle struct {
	TimeTag  TimeTag
	Elements []Packet
}

// Compile-time check to ensure Bundle implements the Packet interface.
var _ Packet = &Bundle{}

/*
NewBundle returns a bundle with immediate time tag.
*/
func NewBundle() *Bundle {
	return &Bundle{TimeTag: TimeTag{Immediate: true}}
}

/*
NewBundleFromData is a convenience factory to decode a bundle from a byte slice.
*/
func NewBundleFromData(data []byte) (*Bundle, error) {
	bundle := NewBundle()
	err := bundle.UnmarshalBinary(data)

	return bundle, err
}

/*
AddPacket adds a packet to the bundle's child elements.
*/
func (bun *Bundle) AddPacket(p Packet) {
	if bun == nil {
		return
	}

	bun.Elements = append(bun.Elements, p)
}

/*
MarshalBinary encodes the Bundle as per the OSC standard.
*/
func (bun Bundle) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)

	buf.Write(bundleString)
	buf.Write(encodeTimeTag(bun.TimeTag))

	// Encode each child element
	for _, e := range bun.Elements {
		encoded, err := e.MarshalBinary()
		if err != nil {
			return nil, err
		}

		count := uint32(len(encoded))
		binary.Write(buf, binary.BigEndian, count)
		buf.Write(encoded)
	}

	bytes := buf.Bytes()
	return bytes, nil
}

/*
UnmarshalBinary attempts to create a new Bundle from an encoded byte slice.
*/
func (bun *Bundle) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	// Check the bundle identifier
	identifier, err := decodeString(buf)
	if err != nil {
		return err
	} else if identifier != "#bundle" {
		return errors.New("Malformed bundle")
	}

	// Read the time tag
	timeTag, err := decodeTimeTag(buf)
	if err != nil {
		return err
	}

	var elements []Packet

	// Read the bundle's contents
	for {
		// Look for a size count
		var count uint32
		err := binary.Read(buf, binary.BigEndian, &count)
		if err == io.EOF {
			// No more bundle data to read, terminate the loop
			break
		} else if err != nil {
			return err
		}

		// Assign a byte array the exact size
		packetData := make([]byte, count)
		n, err := buf.Read(packetData)
		if err != nil {
			return errors.New("Malformed bundle")
		}

		// Ensure that the number of bytes read equals the expected count
		if uint32(n) != count {
			return errors.New("Malformed bundle")
		}

		p, err := decodePacket(packetData)
		if err != nil {
			return err
		}

		elements = append(elements, p)
	}

	bun.TimeTag = timeTag
	bun.Elements = elements

	return nil
}

/*
decodePacket attempts to decode a packet into a Message or a Bundle.
*/
func decodePacket(data []byte) (Packet, error) {
	// Ensure there is data to read, and ensure it is a multiple of 32 bits
	lenData := len(data)
	if lenData <= 0 || lenData%4 != 0 {
		return nil, errors.New("Packet data is not a multiple of 4 bytes")
	}

	var p Packet
	var err error

	firstChar := data[0]
	if firstChar == '/' {
		// The packet is an OSC message
		p, err = NewMessageFromData(data)
		if err != nil {
			return nil, err
		}
	} else if firstChar == '#' {
		// The packet is another bundle
		p, err = NewBundleFromData(data)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Malformed packet")
	}

	return p, nil
}

func (bun Bundle) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString("Bundle: {")
	for _, e := range bun.Elements {
		switch e.(type) {
		case *Bundle:
			buf.WriteString(e.(*Bundle).String())
		case *Message:
			buf.WriteString(e.(*Message).String())
		}
	}
	buf.WriteString("}")

	return buf.String()
}

/*
Equals returns true if bun is equal to other, otherwise false.
*/
func (bun *Bundle) Equals(other *Bundle) bool {
	if &bun == &other {
		return true
	}

	timeTagEq := bun.TimeTag == other.TimeTag

	elementsEq := reflect.DeepEqual(bun.Elements, other.Elements)

	return timeTagEq && elementsEq
}
