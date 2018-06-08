package osc

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

/*
Message represents a single OSC message with address pattern and arguments.
*/
type Message struct {
	Address   string
	Arguments []interface{}
}

/*
NewEmptyMessage returns an OSC message with default values.
*/
func NewEmptyMessage() Message {
	return NewMessage("/")
}

/*
NewMessage creates a new OSC message with an address pattern, and empty arguments.
*/
func NewMessage(address string) Message {
	msg := Message{Address: address}
	return msg
}

/*
NewMessageFromData attempts to create a new Message from an encoded byte slice.
*/
func NewMessageFromData(data []byte) (Message, error) {
	// Create a blank message (note: invalid OSC)
	msg := Message{Address: "", Arguments: nil}

	buf := bytes.NewBuffer(data)

	address, err := decodeString(buf)
	if err != nil {
		return msg, err
	}
	msg.Address = address

	typeTagString, err := decodeString(buf)
	if err != nil {
		return msg, err
	}

	args, err := readArguments(typeTagString, buf)
	if err != nil {
		return msg, err
	}
	msg.Arguments = args

	return msg, nil
}

/*
AddArgument appends a value to the Message's Arguments.
*/
func (msg *Message) AddArgument(arg interface{}) error {
	// If we can get a type tag for the argument, then it is a supported type
	_, err := typeTag(arg)
	if err != nil {
		return fmt.Errorf("Argument type \"%T\" not supported", arg)
	}

	msg.Arguments = append(msg.Arguments, arg)

	return nil
}

/*
TypeTagString generates the type tag string for the Message Arguments.
*/
func (msg *Message) TypeTagString() (string, error) {
	typeTagString := ","

	for _, arg := range msg.Arguments {
		tag, err := typeTag(arg)

		if err != nil {
			return "", err
		}

		typeTagString += tag
	}

	return typeTagString, nil
}

/*
AddressParts returns an ordered array of the individual address parts of the Address of msg.
*/
func (msg *Message) AddressParts() []string {
	parts := strings.Split(msg.Address, "/")
	return parts
}

/*
String implements the fmt.Stringer interface.
*/
func (msg Message) String() string {
	return fmt.Sprintf("Message{Address: %s, Arguments: %v}", msg.Address, msg.Arguments)
}

/*
Bytes encodes the Message as per the OSC standard.
*/
func (msg *Message) Bytes() (data []byte, err error) {
	buf := new(bytes.Buffer)

	buf.Write(encodeString(msg.Address))

	typeTagString, err := msg.TypeTagString()
	if err != nil {
		return nil, err
	}
	buf.Write(encodeString(typeTagString))

	for _, arg := range msg.Arguments {
		argData, err := encodeArgument(arg)
		if err != nil {
			return nil, err
		}
		buf.Write(argData)
	}

	bytes := buf.Bytes()
	return bytes, nil
}

/*
Equals returns true if msg is equal to other, otherwise false.
*/
func (msg *Message) Equals(other *Message) bool {
	if &msg == &other {
		return true
	}

	addressEq := msg.Address == other.Address

	argsEq := reflect.DeepEqual(msg.Arguments, other.Arguments)

	return addressEq && argsEq
}
