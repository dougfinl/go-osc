package osc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

const (
	// Difference in seconds between the Unix epoch (1970) and OSC epoch (1900)
	unixOSCEpochOffset = 2208988800
	// Number of nanoseconds in 1 second
	nanosPerSecond = 1e9
	// The encoded value of an "immediate" time tag
	timeTagImmediate = 0x01
)

/*
TimeTag represents an OSC time tag with an underlying Go time.Time, and an "immediate" flag.
*/
type TimeTag struct {
	time      time.Time
	Immediate bool
}

/*
NewTimeTag returns a TimeTag with the specified Go Time.
*/
func NewTimeTag(t time.Time) TimeTag {
	return TimeTag{time: t, Immediate: false}
}

/*
NewImmediateTimeTag returns a TimeTag representing immediate execution.
*/
func NewImmediateTimeTag() TimeTag {
	return TimeTag{Immediate: true}
}

func (tt TimeTag) String() string {
	var str string

	if tt.Immediate {
		str = "TimeTag: (immediate)"
	} else {
		str = "TimeTag: " + tt.time.String()
	}

	return str
}

/*
typeTag returns the appropriate OSC type tag for a value.
*/
func typeTag(argument interface{}) (string, error) {
	typetag := ""
	var err error

	switch argType := argument.(type) {
	case nil:
		typetag = "N"
	case int32:
		typetag = "i"
	case float32:
		typetag = "f"
	case string:
		typetag = "s"
	case []byte:
		typetag = "b"
	case bool:
		val := argument.(bool)
		if val {
			typetag = "T"
		} else {
			typetag = "F"
		}
	case int64:
		typetag = "h"
	case float64:
		typetag = "d"
	case TimeTag:
		typetag = "t"
	default:
		typetag = ""
		err = fmt.Errorf("Unsupported type: %T", argType)
	}

	return typetag, err
}

/*
encodeString converts an argument to a byte slice.
*/
func encodeArgument(argument interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)

	switch argument.(type) {
	case nil:
		// no bytes are allocated in the argument data
	case int32:
		binary.Write(buf, binary.BigEndian, argument.(int32))
	case float32:
		binary.Write(buf, binary.BigEndian, argument.(float32))
	case string:
		// sequence of non-null ASCII characters followed by a null, followed by 0-3 additional null characters to make
		// the total number of bits a multiple of 32
		buf.Write(encodeString(argument.(string)))
	case []byte:
		// int32 size count, followed by that many 8-bit bytes of arbitrary binary data, followed by 0-3 additional
		// zero bytes to make the total number of bits a multiple of 32
		buf.Write(encodeByteSlice(argument.([]byte)))
	case bool:
		// no bytes are allocated in the argument data
	case int64:
		binary.Write(buf, binary.BigEndian, argument.(int64))
	case float64:
		binary.Write(buf, binary.BigEndian, argument.(float64))
	case TimeTag:
		buf.Write(encodeTimeTag(argument.(TimeTag)))
	default:
		return nil, fmt.Errorf("Unsupported argument type \"%T\"", argument)
	}

	return buf.Bytes(), nil
}

/*
decodeString reads a 32-bit padded OSC string from a byte slice.
*/
func decodeString(buf *bytes.Buffer) (string, error) {
	stringNullTerm, err := buf.ReadString('\x00')

	// Read a null-terminated string
	if err != nil {
		return "", err
	}

	// Trim the null-termination character
	str := strings.Trim(stringNullTerm, "\x00")

	// Calculate how many more null characters we expect to pop (padded to 32 bits)
	stringLength := len(stringNullTerm)
	paddedLength := (stringLength + 3) &^ 0x03

	// Pop the padding, and ensure the values are null
	toPop := paddedLength - stringLength
	for toPop > 0 {
		b, err := buf.ReadByte()
		if b != '\x00' {
			return "", fmt.Errorf("Found a malformed OSC string: %s", err.Error())
		}
		toPop--
	}

	return str, nil
}

/*
readArguments reads a slice of OSC arguments (specific by the typeTagString) from a buffer. If the arguments do not
match the typeTagString, an error is returned.
*/
func readArguments(typeTagString string, buf *bytes.Buffer) ([]interface{}, error) {
	var args []interface{}

	// Ensure the type tag string starts with a comma
	first := typeTagString[:1]
	if first != "," {
		return nil, fmt.Errorf("Malformed type tag string")
	}

	// Iterate over the remaining type tags
	for _, typeTag := range typeTagString[1:] {
		var err error

		switch typeTag {
		case 'T':
			args = append(args, true)
		case 'F':
			args = append(args, true)
		case 'N':
			args = append(args, nil)
		case 'i':
			var val int32
			err = binary.Read(buf, binary.BigEndian, &val)
			args = append(args, val)
		case 'f':
			var val float32
			err = binary.Read(buf, binary.BigEndian, &val)
			args = append(args, val)
		case 's':
			var val string
			val, err = decodeString(buf)
			args = append(args, val)
		case 'b':
			var val []byte
			val, err = decodeByteSlice(buf)
			args = append(args, val)
		case 'h':
			var val int64
			err = binary.Read(buf, binary.BigEndian, &val)
			args = append(args, val)
		case 'd':
			var val float64
			err = binary.Read(buf, binary.BigEndian, &val)
			args = append(args, val)
		case 't':
			var val TimeTag
			val, err = decodeTimeTag(buf)
			args = append(args, val)
		default:
			err = fmt.Errorf("Found unsupported argument type")
		}

		if err != nil {
			return nil, fmt.Errorf("Found malformed argument")
		}
	}

	return args, nil
}

/*
encodeString converts a Go string to a 32-bit padded OSC String.
*/
func encodeString(s string) []byte {
	nullTerminated := []byte(s + string('\x00'))
	return padTo32Bits(nullTerminated)
}

/*
encodeByteSlice converts a Go byte slice to an OSC byte array.
*/
func encodeByteSlice(data []byte) []byte {
	buf := new(bytes.Buffer)
	n := int32(len(data))

	binary.Write(buf, binary.BigEndian, n)
	buf.Write(data)

	paddedBytes := padTo32Bits(buf.Bytes())
	return paddedBytes
}

/*
encodeTimeTag converts a TimeTag to a 64-bit OSC timetag.
*/
func encodeTimeTag(tt TimeTag) []byte {
	var timeTag64 uint64

	if tt.Immediate {
		// If the TimeTag has the "immediate" flag set, ignore the time value
		timeTag64 = timeTagImmediate
	} else {
		// Encode the time with reference to the OSC epoch
		timeOSCSecs := uint64(tt.time.Unix() + unixOSCEpochOffset)
		timeOSCNanos := uint64(tt.time.UnixNano()+unixOSCEpochOffset*nanosPerSecond) - timeOSCSecs*nanosPerSecond

		timeTag64 = timeOSCSecs<<32 | timeOSCNanos&0xFFFFFFFF
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, timeTag64)

	return buf.Bytes()
}

func decodeTimeTag(buf *bytes.Buffer) (TimeTag, error) {
	var timeTag64 uint64

	err := binary.Read(buf, binary.BigEndian, &timeTag64)
	if err != nil {
		return TimeTag{}, err
	}

	var timeTag TimeTag

	if timeTag64 == timeTagImmediate {
		timeTag = NewImmediateTimeTag()
	} else {
		seconds := int64(timeTag64>>32) - unixOSCEpochOffset
		nanoSeconds := int64(timeTag64 & 0xFFFFFFFF)

		t := time.Unix(seconds, nanoSeconds).In(time.UTC)
		timeTag = NewTimeTag(t)
	}

	return timeTag, nil
}

/*
decodeByteSlice reads an OSC byte array into a Go byte slice.
*/
func decodeByteSlice(buf *bytes.Buffer) ([]byte, error) {
	var n int32
	err := binary.Read(buf, binary.BigEndian, &n)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, nil
	}

	// Increase n to the next fourth byte
	nExpected := int((n + 3) &^ 0x03)

	data := make([]byte, nExpected)
	nRead, err := buf.Read(data)
	if err != nil {
		return nil, err
	} else if nRead != nExpected {
		return nil, fmt.Errorf("Didn't read expected number of bytes")
	}

	// Return the slice of the data part of the count
	return data[:n], nil
}

/*
padTo32Bits pads a byte slice to 32 bits by appending nil values.
*/
func padTo32Bits(data []byte) []byte {
	origLength := len(data)

	// Bit-twiddle to find the next multiple of 4 (4 bytes = 32 bits)
	padLength := (origLength + 3) &^ 0x03

	i := padLength - origLength - 1
	for i >= 0 {
		data = append(data, byte(0))
		i--
	}

	return data
}
